package main

import (
	"encoding/json"
	"net/http"

	"github.com/igvaquero18/telegram-notifier/telegram"
)

// SendMessage is a *mux.HandlerFunc that allows us to send notifications with
// an API call
func SendMessage(rw http.ResponseWriter, r *http.Request) {
	message := &telegram.Message{}
	if HandleError(json.NewDecoder(r.Body).Decode(message), rw) {
		return
	}
	if HandleErrorWithStatusCode(telegramClient.SendMessage(message), rw, http.StatusInternalServerError) {
		return
	}
}
