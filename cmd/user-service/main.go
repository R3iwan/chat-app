package main

import (
	"log"
	"net/http"

	"github.com/R3iwan/chat-app/internal/chat"
	"github.com/R3iwan/chat-app/internal/db"
	"github.com/R3iwan/chat-app/internal/message"
	"github.com/R3iwan/chat-app/internal/middleware"
	"github.com/R3iwan/chat-app/internal/user"
	"github.com/R3iwan/chat-app/pkg/config"
	"github.com/R3iwan/chat-app/pkg/logger"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	logger.InitLogger()

	if err != db.ConnectPostgres(cfg) {
		log.Fatalf("could not connect to Postgres: %v", err)
	}
	defer db.ClosePostgres()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/register", user.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", user.LoginHandler).Methods("POST")
	r.Handle("/api/v1/protected", middleware.JWTMiddleware(http.HandlerFunc(protectedHandler))).Methods("GET")
	r.HandleFunc("/ws", chat.WebSocketHandler)
	r.HandleFunc("/ap1/v1/messages", message.GetMessagesHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./frontend"))
	r.PathPrefix("/").Handler(fs)

	log.Printf("User service running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Protected endpoint"))
}
