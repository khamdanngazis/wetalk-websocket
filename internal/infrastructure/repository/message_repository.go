package repository

import (
	"chat-websocket/internal/domain"
	"chat-websocket/internal/domain/entities"
	"sync"
)

type MessageRepositoryImpl struct {
	Messages []entities.Message
	Mutex    sync.Mutex
}

func NewMessageRepository() domain.MessageRepository {
	return &MessageRepositoryImpl{
		Messages: []entities.Message{},
	}
}

func (r *MessageRepositoryImpl) SaveMessage(message entities.Message) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Messages = append(r.Messages, message)
	return nil
}

func (r *MessageRepositoryImpl) GetAllMessages() ([]entities.Message, error) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	return r.Messages, nil
}
