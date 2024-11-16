package controller

import (
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	registerDto := &model.RegisterDTO{}
	if err := c.BodyParser(registerDto); err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	hashedPassword, err := GenerateHashedPassword(registerDto.User.Password)
	user := &model.User{
		Email:    registerDto.User.Email,
		Username: registerDto.User.Username,
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
