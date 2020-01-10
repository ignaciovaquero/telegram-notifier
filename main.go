package main

import (
	"go.uber.org/zap"

	cli "github.com/urfave/cli/v2"
)

var (
	app      *cli.App
	port     int
	chatID   int64
	loglevel bool
	sugar    *zap.SugaredLogger
)

const (
	apiVersion = "/v1"
)

func init() {
	app = &cli.App{
		Name:    "webdding-bot",
		Version: "0.1.0",
		Authors: []*cli.Author{
			{
				Name:  "Ignacio Vaquero Guisasola",
				Email: "ivaqueroguisasola@gmail.com",
			},
		},
		Usage:                "Webdding Telegram Bot",
		UsageText:            "webdding-bot - Webdding Telegram Bot that allows for notifications",
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
				&cli.IntFlag{
					Name:        "port",
					Aliases:     []string{"p"},
					Value:       8081,
					Usage:       "bot listen `PORT`",
					EnvVars:     []string{"WEBDDING_BOT_LISTEN_PORT"},
					DefaultText: "8081",
					Destination: &port,
				},
				&cli.Int64Flag{
					Name:        "chat-id",
					Aliases:     []string{"i"},
					Required:    true,
					Usage:       "Telegram chat `ID` for sending messages",
					Destination: &chatID,
				},
				&cli.BoolFlag{
					Name:        "verbose",
					Aliases:     []string{"v"},
					Value:       false,
					Usage:       "Enable debug logs",
					EnvVars:     []string{"WEBDDING_BOT_VERBOSE"},
					DefaultText: "false",
					Destination: &loglevel,
				},
			},
		},
	}
}

func printError(err error) {
	if err != nil {
		sugar.Fatalf("FATAL: %+v\n", err)
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
	if loglevel {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	zl, err = cfg.Build()

	printError(err)

	sugar = zl.Sugar()
	sugar.Debug("Logger initialization successful")

	return nil
}

func main() {

}
