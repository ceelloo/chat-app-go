package main

import (
	"time"

	"github.com/ceelloo/chat-app-go/internal/store"
	"github.com/ceelloo/chat-app-go/internal/utils"
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

	if _, err := app.store.Users.GetByEmail(c.Context(), payload.Email); err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "User already exist.")
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
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user.")
	}

	expiredAt := time.Now().Add(time.Hour * 24)

	session := store.Session{
		Id:        utils.GenerateId(),
		UserId:    user.Id,
		Token:     utils.GenerateToken(32),
		CsrfToken: utils.GenerateToken(32),
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
		"status":     "success",
		"data":       user,
		"token":      createdSession.Token,
		"csrf_token": createdSession.CsrfToken,
		"message":    "User created successfully",
	})
}

type LoginUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) loginUserHandler(c *fiber.Ctx) error {
	var payload LoginUserPayload

	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid payload")
	}

	user, err := app.store.Users.GetByEmail(c.Context(), payload.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)
	expiredAt := time.Now().Add(time.Hour * 24)

	session := store.Session{
		Id:        utils.GenerateId(),
		UserId:    user.Id,
		Token:     sessionToken,
		CsrfToken: csrfToken,
		CreatedAt: time.Now().Format(time.RFC3339),
		ExpiresAt: expiredAt.String(),
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

	return c.JSON(fiber.Map{
		"message":    "login successful",
		"csrf_token": csrfToken,
	})
}

func (app *application) authenticateUserHandler(c *fiber.Ctx) error {

	user, err := app.store.Users.GetById(c.Context(), c.Locals("userId").(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"data":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}
