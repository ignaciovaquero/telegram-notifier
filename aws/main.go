package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

const (
	verboseEnv string = "NOTIFIER_BOT_VERBOSE"
	tokenEnv   string = "NOTIFIER_BOT_TOKEN"
)

var (
	errorMsg            string
	internalServerError *events.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
	}
)

func getOrElse(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	v := getOrElse(verboseEnv, "true")
	verbose, err := strconv.ParseBool(v)
	if err != nil {
		fmt.Printf("Incorrect value for %s: %s. Defaulting to false", verboseEnv, v)
		verbose = false
	}

	var zl *zap.Logger
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
		errorMsg = fmt.Sprintf("Error when initializing logger: %v", err)
		internalServerError.Body = errorMsg
		return internalServerError, fmt.Errorf(errorMsg)
	}

	sugar := zl.Sugar()
	sugar.Debug("Logger initialization successful")

	telegramClient, err := telegram.NewTelegram(os.Getenv(tokenEnv), sugar)
	if err != nil {
		errorMsg = fmt.Sprintf("Error when initializing the Telegram Client: %s", err.Error())
		internalServerError.Body = errorMsg
		return internalServerError, fmt.Errorf(errorMsg)
	}

}

func main() {
	lambda.Start(Handler)
}
