package user

import (
	"context"
	"encoding/json"
	"fmt"
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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
		log.Printf("Error fetching password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Hashed password:", hashedPassword)

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := getUserID(req.Username)
	if err != nil {
		http.Error(w, "Internal server error: could not get user ID", http.StatusInternalServerError)
		return
	}

	accessToken, err := generateJWT(req.Username, userID, 15*time.Minute)
	if err != nil {
		http.Error(w, "Internal server error: could not generate JWT", http.StatusInternalServerError)
		return
	}

	refreshToken, err := generateJWT(req.Username, userID, 7*24*time.Hour)
	if err != nil {
		http.Error(w, "Internal server error: could not generate JWT", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Printf("User %s logged in", req.Username)

}

func getUserID(username string) (int, error) {
	var userID int
	err := db.DB.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		log.Printf("Error getting user ID for username %s: %v", username, err)
		return 0, err
	}
	return userID, nil
}

func getUserHashPassword(username string) (string, error) {
	var hashedPassword string
	err := db.DB.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Error getting user password for username %s: %v", username, err)
		return hashedPassword, nil
	}

	return hashedPassword, nil
}

func generateJWT(username string, userID int, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"userID":   userID,
		"exp":      time.Now().Add(time.Hour * 3).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractUserIDFromJWT(tokenStr string) (float64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		return 0, fmt.Errorf("userID not found or invalid type")
	}

	return userID, nil
}
