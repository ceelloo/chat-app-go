package main

import "github.com/gofiber/fiber/v2"

func (app *application) Authorize(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	csrfToken := c.Get("X-CSRF-TOKEN")

	if sessionToken == "" || csrfToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	session, err := app.store.Sessions.Get(c.Context(), sessionToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	if session.CsrfToken != csrfToken {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid csrf token",
		})
	}

	c.Locals("userId", session.UserId)

	return c.Next()
}

func (app *application) ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal server error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	} else if err != nil {
		msg = err.Error()
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": msg,
		"code":    code,
	})
}
