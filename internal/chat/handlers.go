package chat

import (
	"log"
	"net/http"

	"github.com/R3iwan/chat-app/internal/message"
	"github.com/R3iwan/chat-app/internal/middleware"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type User struct {
	Conn *websocket.Conn
	Send chan []byte
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	claims, err := middleware.ValidateJWT(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Could not read message", err)
			return
		}
		log.Printf("Message received: %s", string(msg))

		if err = conn.WriteMessage(messageType, msg); err != nil {
			log.Println("Could not write message", err)
			return
		}

		err = message.SaveMessage(int(userID), 2, string(msg))
		if err != nil {
			log.Println("Could not save message", err)
			return
		}

		log.Printf("Connection ended for %d", int(userID))
	}
}
