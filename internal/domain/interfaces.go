package domain

import "chat-websocket/internal/domain/entities"

type MessageRepository interface {
	SaveMessage(message entities.Message) error
	GetAllMessages() ([]entities.Message, error)
}

type WebSocketHandler interface {
	HandleConnections()
	BroadcastMessages(paths []string)
}

type Producer interface {
	SendMessage(topic string, key, message []byte) error
}
