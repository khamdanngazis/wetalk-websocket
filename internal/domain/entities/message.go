package entities

import (
	"time"
)

type WSMessage struct {
	ID           string         `json:"id"`
	Sender       string         `json:"sender_id"`
	ChatRoom     string         `json:"chat_room_id"`
	Receiver     string         `json:"receiver_id"`
	Participants []Participants `json:"participants"`
	Status       int            `json:"status"`
	Content      string         `json:"content"`
	Timestamp    string         `json:"timestamp"`
}

type Participants struct {
	UserID     string `json:"user_id"`
	SocketPath string `json:"socket_path"`
}

type Message struct {
	ID         string    `json:"id"`
	ChatRoomID string    `json:"chat_room_id"` // Referensi ke ChatRoom
	SenderID   string    `json:"sender_id"`    // ID pengirim pesan
	Content    string    `json:"content"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type MessageStatus struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	MessageID  string `json:"message_id"`
	ReceiverID string `json:"receiver_id"`
	Status     int    `json:"status"`
}

const (
	StatusSend      = 1
	StatusDelivered = 2
	StatusRead      = 3
)
