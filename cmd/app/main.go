package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/princecee/go_chat/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// set the configs and dependency injections
	app.StartApp()
}
