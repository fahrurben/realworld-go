package controller

import (
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/pkg/validator"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Register(c *fiber.Ctx) error {
	registerDto := &model.RegisterDTO{}
	if err := c.BodyParser(registerDto); err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	// Create a new validator for a Book model.
	validate := validator.NewValidator()
	if err := validate.Struct(registerDto.User); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	hashedPassword, err := GenerateHashedPassword(registerDto.User.Password)
	user := &model.User{
		Email:    strings.ToLower(registerDto.User.Email),
		Username: strings.ToLower(registerDto.User.Username),
		Password: hashedPassword,
	}
	userRepo := repository.NewUserRepo(database.GetDB())

	if exists, _ := userRepo.Exists(user.Username, user.Email); exists == true {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, "Email or Username already exists")
	}

	id, err := userRepo.Create(user)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	createdUser, err := userRepo.Get(id)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	token, err := GenerateAccessToken(createdUser)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(fiber.Map{"user": fiber.Map{
		"email":    createdUser.Email,
		"token":    token,
		"username": createdUser.Username,
		"bio":      createdUser.Bio,
		"image":    createdUser.Image,
	}})
}
