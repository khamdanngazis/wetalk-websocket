package websocket

import (
	"chat-websocket/internal/application/usecases"
	"chat-websocket/internal/config"
	"chat-websocket/internal/domain"
	"chat-websocket/internal/domain/entities"
	"chat-websocket/package/helper"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketHandlerImpl struct {
	Upgrader     websocket.Upgrader
	ClientGroups map[string]map[string]*websocket.Conn
	Broadcast    chan entities.WSMessage
	Mutex        sync.Mutex
	ChatUsecase  *usecases.ChatUsecase
}

func NewWebSocketHandler(chatUsecase *usecases.ChatUsecase) domain.WebSocketHandler {
	return &WebSocketHandlerImpl{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		ClientGroups: make(map[string]map[string]*websocket.Conn),
		Broadcast:    make(chan entities.WSMessage),
		ChatUsecase:  chatUsecase,
	}
}

func (h *WebSocketHandlerImpl) HandleConnections() {
	go func() {
		for {
			message := <-h.Broadcast
			h.Mutex.Lock()
			message.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			for _, v := range message.Participants {
				// Check if the client exists in the group
				client, ok := h.ClientGroups[v.SocketPath][v.UserID]
				if !ok {
					log.Printf("Client not found for SocketPath: %s, UserID: %s", v.SocketPath, v.UserID)
					continue // Skip to the next participant
				}

				// Attempt to send the message
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("Error writing message to client: %v", err)
					client.Close()
					delete(h.ClientGroups[v.SocketPath], v.UserID) // Remove client from the group
				}
			}
			h.Mutex.Unlock()
		}
	}()
}

func (h *WebSocketHandlerImpl) BroadcastMessages(paths []string) {

	for _, path := range paths {
		fmt.Println("Broadcasting to path:", path)
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			conn, err := h.Upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
			queryParams := r.URL.Query()
			tokenString := queryParams.Get("token")

			if tokenString == "" {
				log.Println("Error: Missing Authorization")
				conn.Close() // Close the connection immediately if ID is missing
				return
			}

			// Validate the token
			user, err := helper.ValidateToken(tokenString)
			if err != nil {
				log.Println("Error: Not Valid Authorization")
				conn.Close() // Close the connection immediately if ID is missing
				return
			}

			if user.UserID == "" {
				log.Println("Error: User ID not found")
				conn.Close() // Close the connection immediately if ID is missing
				return
			}

			h.Mutex.Lock()
			if _, exists := h.ClientGroups[path]; !exists {
				h.ClientGroups[path] = make(map[string]*websocket.Conn)
			}
			if existingConn, exists := h.ClientGroups[path][user.UserID]; exists {
				fmt.Printf("Client %s already connected, closing existing connection\n", user.UserID)
				if existingConn != nil {
					err := existingConn.Close()
					if err != nil {
						fmt.Printf("Error closing existing connection for user %s: %v\n", user.UserID, err)
					}
				}
				delete(h.ClientGroups[path], user.UserID)
			}
			h.ClientGroups[path][user.UserID] = conn
			fmt.Printf("Client connected: UserID=%s, Username=%s\n", user.UserID, user.Username)
			fmt.Printf("Total clients for path '%s': %d\n", path, len(h.ClientGroups[path]))
			h.Mutex.Unlock()

			defer func() {
				h.Mutex.Lock()
				delete(h.ClientGroups[path], user.UserID)
				h.Mutex.Unlock()
				conn.Close()
			}()

			for {
				var message entities.WSMessage
				err := conn.ReadJSON(&message)
				if err != nil {
					log.Printf("error: %v", err)
					break
				}

				if message.Status == 0 {
					message.Status = entities.StatusSend
				} else {

				}
				fmt.Println("receive message ", message)
				h.Broadcast <- message
				if err := h.ChatUsecase.BroadcastMessage(config.GetEnv("KAFKA_TOPIC", "chat"), message); err != nil {
					fmt.Println("Failed to broadcast message:", err)
				}
			}
		})
	}
}
