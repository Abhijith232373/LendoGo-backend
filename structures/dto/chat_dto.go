package dto

import "time"

type IncomingMessage struct {
	Text        string `json:"text"`
	IsFromAdmin bool   `json:"is_from_admin"`
	ReceiverID  string `json:"receiver_id"` 
}

type OutgoingMessage struct {
	SenderID    string    `json:"sender_id"`
	ReceiverID  string    `json:"receiver_id"`
	IsFromAdmin bool      `json:"is_from_admin"`
	Text        string    `json:"text"`
	Timestamp   time.Time `json:"timestamp"`
}