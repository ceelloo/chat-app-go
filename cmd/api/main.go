package main

import (
	"log"

	"github.com/ceelloo/chat-go/internal/database"
	"github.com/ceelloo/chat-go/internal/env"
	"github.com/ceelloo/chat-go/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":3000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "local.db"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := database.New(cfg.db.addr, cfg.db.maxIdleTime, cfg.db.maxOpenConns, cfg.db.maxIdleConns)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if cfg.env == "development" {
		if err := database.LoadMigration("./cmd/migrate/migrations", db); err != nil {
			log.Fatal(err)
		}
	}

	store := store.NewStorage(db)

	app := application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	app.serve(mux)
}
