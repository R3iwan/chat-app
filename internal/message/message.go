package message

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/R3iwan/chat-app/internal/db"
	"github.com/R3iwan/chat-app/internal/user"
)

type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
}

func SaveMessage(senderID, receiverID int, content string) error {
	_, err := db.DB.Exec(context.Background(),
		"INSERT INTO messages (sender_id, receiver_id, content, sent_at) VALUES ($1, $2, $3, $4)",
		senderID, receiverID, content, time.Now(),
	)
	return err
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	senderID := r.URL.Query().Get("sender_id")
	receiverID := r.URL.Query().Get("receiver_id")

	rows, err := db.DB.Query(context.Background(),
		`SELECT sender_id, receiver_id, content, sent_at FROM messages
	WHERE (sender_id = $1) AND (receiver_id = $2) OR (sender_id = $2) AND (receiver_id = $1)
	ORDER BY sent_at`, senderID, receiverID)

	if err != nil {
		http.Error(w, "Internal Response Error", http.StatusInternalServerError)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.SentAt); err != nil {
			http.Error(w, "Internal Response Error", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := user.ExtractUserIDFromJWT(tokenString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	msg.SenderID = int(userID)
	msg.SentAt = time.Now()

	query := `INSERT INTO messages (sender_id, receiver_id, content, sent_at) VALUES ($1, $2, $3, $4)`
	_, err = db.DB.Exec(context.Background(), query, msg.SenderID, msg.ReceiverID, msg.Content, msg.SentAt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
