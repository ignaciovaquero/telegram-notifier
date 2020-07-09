package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/igvaquero18/telegram-notifier/telegram"
	"go.uber.org/zap"
)

const (
	tokenEnv   string = "NOTIFIER_BOT_TOKEN"
	verboseEnv string = "NOTIFIER_BOT_VERBOSE"
)

var (
	token               string = os.Getenv(tokenEnv)
	v                   string = getOrElse(verboseEnv, "true")
	sugar               *zap.SugaredLogger
	errorMsg            string
	telegramClient      *telegram.Client
	verbose             bool
	internalServerError *events.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
	}
	badRequest *events.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadGateway,
	}
	okResponse *events.APIGatewayProxyResponse = &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}
)

func getOrElse(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func init() {
	var err error

	verbose, err = strconv.ParseBool(v)
	if err != nil {
		log.Printf("Incorrect value for %s: %s. Defaulting to false", verboseEnv, v)
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
		log.Fatalf("Error when initializing logger: %v", err)
	}

	sugar = zl.Sugar()
	sugar.Debug("Logger initialization successful")

	telegramClient, err = telegram.NewClient(token, sugar)
	if err != nil {
		log.Fatalf("Error when initializing the Telegram Client: %s", err.Error())
	}
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	message := &telegram.Message{}

	if err := json.Unmarshal([]byte(request.Body), message); err != nil {
		badRequest.Body = err.Error()
		return badRequest, err
	}

	if err := telegramClient.SendMessage(message); err != nil {
		internalServerError.Body = err.Error()
		return internalServerError, err
	}
	return okResponse, nil
}

func main() {
	lambda.Start(Handler)
}
