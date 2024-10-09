package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/R3iwan/chat-app/internal/db"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(viper.GetString("JWT_SECRET"))

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Bad request: missing fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := getUserHashPassword(req.Username)
	if err != nil {
		http.Error(w, "Internal server error: could not hash password", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := genereteJWT(req.Username)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Printf("User %s logged in", req.Username)

}

func getUserHashPassword(username string) (string, error) {
	var hashedPassword string
	err := db.DB.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Error getting user password for username %s: %v", username, err)
		return "", err
	}

	return hashedPassword, nil
}

func genereteJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
