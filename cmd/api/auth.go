package main

import (
	"time"

	"github.com/ceelloo/chat-go/internal/store"
	"github.com/ceelloo/chat-go/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type RegisterUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) registerUserHandler(c *fiber.Ctx) error {
	var payload RegisterUserPayload

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	user := &store.User{
		Id:    utils.GenerateId(),
		Name:  payload.Name,
		Email: payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		return err
	}

	if err := app.store.Users.Create(c.Context(), *user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	expiredAt := time.Now().Add(time.Hour * 24)

	session := store.Session{
		Id:        utils.GenerateId(),
		UserId:    user.Id,
		Token:     utils.GenerateToken(32),
		CsrfToken: utils.GenerateToken(16),
		ExpiresAt: expiredAt.String(),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	createdSession, err := app.store.Sessions.Create(c.Context(), session)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create session")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    createdSession.Token,
		HTTPOnly: true,
		Expires:  expiredAt,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "csrf_token",
		Value:    createdSession.CsrfToken,
		HTTPOnly: false,
		Expires:  expiredAt,
	})

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   user,
		"token":  createdSession.Token,
		"csrf_token": createdSession.CsrfToken,
		"message": "User created successfully",
	})
}
