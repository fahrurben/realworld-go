package controller

import "github.com/gofiber/fiber/v2"

func CreateErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"errors": fiber.Map{
			"body": []string{message},
		},
	})
}
