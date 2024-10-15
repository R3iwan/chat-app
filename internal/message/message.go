package message

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/R3iwan/chat-app/internal/db"
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
