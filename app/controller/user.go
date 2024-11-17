package controller

import (
	"fmt"
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/golang-jwt/jwt/v4"
)
import "github.com/gofiber/fiber/v2"

func GetCurrentUser(c *fiber.Ctx) error {
	fmt.Println(c.Locals("user"))
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := int64(claims["user_id"].(float64))

	userRepo := repository.NewUserRepo(database.GetDB())

	user, err := userRepo.Get(int64(user_id))
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(fiber.Map{"user": fiber.Map{
		"email":    user.Email,
		"token":    "",
		"username": user.Username,
		"bio":      user.Bio,
		"image":    user.Image,
	}})
}
