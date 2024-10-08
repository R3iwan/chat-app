package main

import (
	"log"
	"net/http"

	"github.com/R3iwan/chat-app/internal/db"
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
	r.HandleFunc("/api/v1/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", LoginHandler).Methods("POST")

	log.Printf("User service running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("register"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("login"))
}
