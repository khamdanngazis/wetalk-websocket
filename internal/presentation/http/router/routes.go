package router

import (
	"chat-websocket/internal/application/usecases"
	"chat-websocket/internal/config"
	"chat-websocket/internal/infrastructure/kafka"
	"chat-websocket/internal/infrastructure/websocket"
	"net/http"
)

func SetupRoutes(paths []string) {
	kafkaProducer := kafka.NewKafkaService([]string{config.GetEnv("KAFKA_HOST", "localhost:9092")})
	chatUsecase := usecases.NewChatUsecase(kafkaProducer)
	handler := websocket.NewWebSocketHandler(chatUsecase)
	handler.HandleConnections()
	//paths := []string{"group_1", "group_2", "group_3"}
	handler.BroadcastMessages(paths)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
}
