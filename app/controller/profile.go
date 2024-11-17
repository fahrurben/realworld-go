package controller

import (
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetProfile(c *fiber.Ctx) error {
	var user_id int64 = 0
	if token := c.Locals("user"); token != nil {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		user_id = int64(claims["user_id"].(float64))
	}

	username := c.Params("username")
	userRepo := repository.NewUserRepo(database.GetDB())

	user, err := userRepo.GetByUsername(username)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	isFollowing := false
	if user_id > 0 {
		isFollowing, _ = userRepo.IsFollowing(user_id, user.ID)
	}

	return c.JSON(fiber.Map{"profile": fiber.Map{
		"username":  user.Username,
		"bio":       user.Bio,
		"image":     user.Image,
		"following": isFollowing,
	}})
}