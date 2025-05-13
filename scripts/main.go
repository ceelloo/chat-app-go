package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ceelloo/chat-app-go/internal/database"
	"github.com/ceelloo/chat-app-go/internal/env"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify a command to run")
	}

	command := os.Args[1]

	switch command {
		case "db-reset":
			fmt.Println("Resetting database")

			db, err := database.New(
				env.GetString("DB_ADDR", "local.db"),
				env.GetString("DB_MAX_IDLE_TIME", "15m"),
				env.GetInt("DB_MAX_OPEN_CONNS", 25),
				env.GetInt("DB_MAX_IDLE_CONNS", 25),
			)
			if err != nil {
				log.Fatal(err)
			}

			defer db.Close()

			if err := database.LoadMigration("./cmd/migrate/migrations", db); err != nil {
				log.Fatal(err)
			}
	}
}