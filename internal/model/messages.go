package model

import "time"

type Messages struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Text       string    `json:"text"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
