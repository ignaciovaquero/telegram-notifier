package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fmt"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/igvaquero18/telegram-notifier/controllers"
	cli "github.com/urfave/cli/v2"
)

var (
	app            *cli.App
	wait           time.Duration
	address, token string
	port           int
	verbose        bool
	sugar          *zap.SugaredLogger
)

const (
	apiVersion = "/v1"
)

func init() {
	app = &cli.App{
		Name:    "telegram-notifier",
		Version: "0.1.0",
		Authors: []*cli.Author{
			{
				Name:  "Ignacio Vaquero Guisasola",
				Email: "ivaqueroguisasola@gmail.com",
			},
		},
		Usage:                "Telegram Notifier Bot",
		UsageText:            "telegram-notifier - Telegram Bot that allows for notifications",
		EnableBashCompletion: true,
	}
	app.Commands = []*cli.Command{
		{
			Name:        "start",
			Usage:       "starts the server",
			UsageText:   "starts the server and listens to connections",
			Description: "this command starts the server. Parameters can be specified through the use of flags, for DB connection and others",
			Action:      run,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "address",
					Aliases:     []string{"a"},
					Value:       "0.0.0.0",
					Usage:       "bot listen `ADDRESS`",
					EnvVars:     []string{"NOTIFIER_BOT_LISTEN_ADDRESS"},
					DefaultText: "0.0.0.0",
					Destination: &address,
				},
				&cli.IntFlag{
					Name:        "port",
					Aliases:     []string{"p"},
					Value:       8081,
					Usage:       "bot listen `PORT`",
					EnvVars:     []string{"NOTIFIER_BOT_LISTEN_PORT"},
					DefaultText: "8081",
					Destination: &port,
				},
				&cli.StringFlag{
					Name:        "token",
					Aliases:     []string{"k"},
					Required:    true,
					Usage:       "the Telegram `TOKEN` to talk to the API",
					EnvVars:     []string{"NOTIFIER_BOT_TOKEN"},
					Destination: &token,
				},
				&cli.DurationFlag{
					Name:        "graceful-timeout",
					Aliases:     []string{"timeout", "t"},
					Value:       time.Second * 15,
					Usage:       "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m",
					EnvVars:     []string{"NOTIFIER_BOT_TIMEOUT"},
					DefaultText: "15s",
					Destination: &wait,
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Aliases:     []string{"v"},
					Value:       false,
					Usage:       "enable debug logs",
					EnvVars:     []string{"NOTIFIER_BOT_VERBOSE"},
					DefaultText: "false",
					Destination: &verbose,
				},
			},
		},
	}
}

func run(clictx *cli.Context) error {
	// Building the logger
	var zl *zap.Logger
	var err error
	cfg := zap.Config{
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if verbose {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	zl, err = cfg.Build()

	if err != nil {
		return fmt.Errorf("error when initializing logger: %v", err)
	}

	sugar = zl.Sugar()
	sugar.Debug("Logger initialization successful")

	telegram, err := controllers.NewTelegram(token, sugar)
	if err != nil {
		return fmt.Errorf("error when initializing Telegram client: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc(fmt.Sprintf("%s/notifications", apiVersion), telegram.SendNotification).Methods(http.MethodPost)
	sugar.Debug("Router setup complete")

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", address, port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	e := make(chan error, 1)

	go func() {
		sugar.Infof("Start listening on %s:%d", address, port)
		if err := srv.ListenAndServe(); err != nil {
			e <- err
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case er := <-e:
		return fmt.Errorf("error when listening at %s:%d -> %v", address, port, er)
	case <-c:
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()

		srv.Shutdown(ctx)

		sugar.Info("shutting down")
		break
	}

	return nil
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
