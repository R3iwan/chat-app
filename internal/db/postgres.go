package db

import (
	"context"
	"log"
	"time"

	"github.com/R3iwan/chat-app/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectPostgres(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	connConfig, err := pgxpool.ParseConfig(cfg.PostgresURL)
	if err != nil {
		return err
	}
	db, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return err
	}

	DB = db
	log.Println("Connected to Postgres")
	return nil
}

func ClosePostgres() {
	DB.Close()
	log.Println("Closed Postgres connection")
}
