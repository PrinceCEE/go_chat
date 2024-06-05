package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/princecee/go_chat/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("DSN")
	pool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	app.StartApp(pool)
}
