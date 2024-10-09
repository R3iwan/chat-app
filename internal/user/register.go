package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/R3iwan/chat-app/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.Email == "" {
		http.Error(w, "Bad request: missing fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal server error: could not hash password", http.StatusInternalServerError)
		log.Printf("Error hashing password: %v", err)
		return
	}

	if err := createUser(req.Username, req.Email, hashedPassword); err != nil {
		log.Printf("Error inserting user: %v", err)
		http.Error(w, "User already exists or server error", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func hashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPassword), nil
}

func createUser(username, email, passwordHash string) error {
	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			log.Printf("Transaction failed, rolling back: %v", err)
			tx.Rollback(context.Background())
		} else {
			log.Printf("Committing transaction")
			tx.Commit(context.Background())
		}
	}()

	query := `
        INSERT INTO users (username, email, password_hash, created_at)
        VALUES ($1, $2, $3, $4)
    `
	log.Printf("Attempting to insert user: username=%s, email=%s", username, email)

	_, err = tx.Exec(context.Background(), query, username, email, passwordHash, time.Now())
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return err
	}

	log.Printf("User inserted successfully")
	return nil
}
