package usecases

import (
	"chat-websocket/internal/domain"
	"chat-websocket/internal/domain/entities"
	"encoding/json"
	"time"
)

type ChatUsecase struct {
	KafkaProducer domain.Producer
}

func NewChatUsecase(kafkaProducer domain.Producer) *ChatUsecase {
	return &ChatUsecase{KafkaProducer: kafkaProducer}
}

func (u *ChatUsecase) BroadcastMessage(topic string, message entities.WSMessage) error {

	if message.Status == entities.StatusSend {
		parsedTime, err := time.Parse(time.RFC3339, message.Timestamp)
		if err != nil {
			return err
		}
		sendMessage := entities.Message{
			ID:         message.ID,
			SenderID:   message.Sender,
			ChatRoomID: message.ChatRoom,
			Content:    message.Content,
			Status:     message.Status,
			CreatedAt:  parsedTime,
		}
		messageBytes, err := json.Marshal(sendMessage)
		if err != nil {
			return err
		}
		return u.KafkaProducer.SendMessage(topic, []byte("message"), messageBytes)
	} else {
		updateMessage := entities.MessageStatus{
			MessageID:  message.ID,
			ReceiverID: message.Receiver,
			Status:     message.Status,
		}
		messageBytes, err := json.Marshal(updateMessage)
		if err != nil {
			return err
		}
		return u.KafkaProducer.SendMessage(topic, []byte("update_status"), messageBytes)
	}
}
