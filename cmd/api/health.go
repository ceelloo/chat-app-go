package main

import "github.com/gofiber/fiber/v2"

func (app *application) healthCheckHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
} 