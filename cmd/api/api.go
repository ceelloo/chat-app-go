package main

import (
	"time"

	"github.com/ceelloo/chat-app-go/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *fiber.App {
	r := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	})

	// Middlewares
	r.Use(requestid.New())
	r.Use(logger.New())
	r.Use(recover.New())

	// Routes
	auth := r.Group("/authentication")
	auth.Post("/register", app.registerUserHandler)

	r.Get("/health", app.healthCheckHandler)

	return r
}

func (app *application) serve(fb *fiber.App) error {
	return fb.Listen(app.config.addr)
}
