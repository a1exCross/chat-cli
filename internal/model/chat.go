package model

import "time"

type Message struct {
	From      string
	Text      string
	CreatedAt time.Time
}

type SendMessageParams struct {
	Message
	ChatID int64
}

type Chat struct {
	ID        int64
	Usernames []string
}
