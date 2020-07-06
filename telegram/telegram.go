package telegram

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Client A client for the Telegram Bot
type Client struct {
	Logger
	*tgbotapi.BotAPI
}

// NewTelegram returns a pointer to a Client object
func NewClient(token string, log Logger) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return &Client{
		Logger: log,
		BotAPI: bot,
	}, nil
}

// NewDefaultClient returns a pointer to a Client object, with default
// logging capabilities.
func NewDefaultClient(token string) (*Client, error) {
	d := &defaultLogger{
		log.New(os.Stdout, "logger: ", log.LstdFlags),
	}
	return NewClient(token, d)
}

// SendMessage sends a Telegram message to a set of chats
func (c *Client) SendMessage(message *Message) error {
	return c.SendNotification(message.Title, message.Message, message.Chats)
}

// SendNotification sends a notification to a set of chats
func (c *Client) SendNotification(title, body string, chats []int64) error {
	var wg sync.WaitGroup
	errors := []string{}

	wg.Add(len(chats))

	c.Debugf("Sending notifications to %d chats", len(chats))
	for _, chat := range chats {
		go func(ch int64) {
			defer wg.Done()
			var msgconf tgbotapi.MessageConfig
			var msg string

			if title != "" {
				msg = fmt.Sprintf("*%s*\n\n%s", title, body)
				msgconf = tgbotapi.NewMessage(ch, msg)
			} else {
				msg = fmt.Sprintf("%s", body)
				msgconf = tgbotapi.NewMessage(ch, msg)
			}

			msgconf.ParseMode = "Markdown"

			c.Debugf("Sending message %s to chat with ID %d", msg, ch)
			_, err := c.Send(msgconf)
			if err != nil {
				c.Debugf("Error when sending message: %s", err.Error())
				errors = append(errors, err.Error())
			}
		}(chat)
	}

	wg.Wait()

	if len(errors) > 0 {
		c.Logger.Debug("Returning errors back to the caller")
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}
