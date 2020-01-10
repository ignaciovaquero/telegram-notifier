package models

// Message a message to send
type Message struct {
	Title   string  `json:"title,omitempty"`
	Message string  `json:"message"`
	Chats   []int64 `json:"chats"`
}
