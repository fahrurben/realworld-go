package controller

import (
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/gofiber/fiber/v2"
)

func ListTags(c *fiber.Ctx) error {
	tagRepo := repository.NewTagRepo(database.GetDB())
	tags, err := tagRepo.List()

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(fiber.Map{"tags": tags})
}
