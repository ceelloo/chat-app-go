package main

import (
	"time"

	"github.com/ceelloo/chat-app-go/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	init := fiber.New(fiber.Config{
		ErrorHandler: app.ErrorHandler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	})

	r := init.Group("/api") 

	// Middlewares
	r.Use(requestid.New())
	r.Use(logger.New())
	r.Use(recover.New())

	r.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // React/Vite dev server
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-CSRF-Token",
		AllowCredentials: true, // Jika kamu pakai cookie
	}))

	// Routes
	auth := r.Group("/authentication")
	auth.Post("/register", app.registerUserHandler)
	auth.Post("/login", app.loginUserHandler)
	auth.Get("/me", app.Authorize, app.authenticateUserHandler)

	r.Get("/health", app.healthCheckHandler)

	return init
}

func (app *application) serve(fb *fiber.App) error {
	return fb.Listen(app.config.addr)
}
