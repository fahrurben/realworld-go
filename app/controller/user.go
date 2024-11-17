package controller

import (
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/pkg/validator"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)
import "github.com/gofiber/fiber/v2"

func GetCurrentUser(c *fiber.Ctx) error {
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

func UpdateUser(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := int64(claims["user_id"].(float64))

	userRepo := repository.NewUserRepo(database.GetDB())

	currentUser, err := userRepo.Get(int64(user_id))
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	updateDto := &model.UpdateDto{}
	if err := c.BodyParser(updateDto); err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	// Create a new validator for a Book model.
	validate := validator.NewValidator()
	if err := validate.Struct(updateDto.User); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	userDto := updateDto.User
	if userDto.Username != "" {
		currentUser.Username = strings.ToLower(userDto.Username)
	}

	if userDto.Password != "" {
		currentUser.Password, _ = GenerateHashedPassword(userDto.Password)
	}

	if userDto.Bio != "" {
		currentUser.Bio = &userDto.Bio
	}

	if userDto.Image != "" {
		currentUser.Image = &userDto.Image
	}

	err = userRepo.Update(int64(currentUser.ID), currentUser)

	return c.JSON(fiber.Map{"user": fiber.Map{
		"email":    currentUser.Email,
		"token":    "",
		"username": currentUser.Username,
		"bio":      currentUser.Bio,
		"image":    currentUser.Image,
	}})
}
