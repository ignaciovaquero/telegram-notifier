package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/igvaquero18/webdding-bot/models"
	"github.com/igvaquero18/webdding-bot/utils"
)

// Telegram A Telegram bot with channel
type Telegram struct {
	Logger
	*tgbotapi.BotAPI
}

// NewTelegram returns a pointer to a Telegram object
func NewTelegram(token string, log Logger) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Telegram{
		Logger: log,
		BotAPI: bot,
	}, nil
}

// NewDefaultTelegram returns a pointer to a Telegram object, with default
// logging capabilities.
func NewDefaultTelegram(token string) (*Telegram, error) {
	d := &defaultLogger{
		log.New(os.Stdout, "logger: ", log.LstdFlags),
	}
	return NewTelegram(token, d)
}

// LoggingMiddleware The logging middleware allows to log any request whenever the log level is equal to Debug.
func (t *Telegram) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		t.Debugf("%s %s", r.Method, r.URL.EscapedPath())
		next.ServeHTTP(rw, r)
	})
}

// SendNotification sends a notification to a set of chats
func (t *Telegram) SendNotification(rw http.ResponseWriter, r *http.Request) {
	message := &models.Message{}
	if utils.HandleError(json.NewDecoder(r.Body).Decode(message), rw) {
		return
	}

	var wg sync.WaitGroup
	errors := []string{}

	wg.Add(len(message.Chats))

	t.Debugf("Sending notifications to %d chats", len(message.Chats))
	for _, chat := range message.Chats {
		go func(ch int64) {
			defer wg.Done()
			var msgconf tgbotapi.MessageConfig
			var msg string

			if message.Title != "" {
				msg = fmt.Sprintf("`%s`\n\n%s", message.Title, message.Message)
				msgconf = tgbotapi.NewMessage(ch, msg)
			} else {
				msg = fmt.Sprintf("%s", message.Message)
				msgconf = tgbotapi.NewMessage(ch, msg)
			}

			msgconf.ParseMode = "Markdown"

			t.Debugf("Sending message %s to chat with ID %d", msg, ch)
			_, err := t.Send(msgconf)
			if err != nil {
				t.Debugf("Error when sending message: %s", err.Error())
				errors = append(errors, err.Error())
			}
		}(chat)
	}

	wg.Wait()

	if len(errors) > 0 {
		t.Logger.Debug("Returning errors back to the caller")
		rw.WriteHeader(http.StatusInternalServerError)
		_, err := io.WriteString(rw, strings.Join(errors, "\n"))
		utils.HandleErrorWithStatusCode(err, rw, http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)
}
