// handlers/status_handler.go
package handlers

import "github.com/gofiber/fiber/v2"

func StatusHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "running",
		"version": "1.0.0",
	})
}
