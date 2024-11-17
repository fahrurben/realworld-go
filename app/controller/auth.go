package controller

import (
	"database/sql"
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/pkg/config"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	loginDto := &model.LoginDTO{}
	if err := c.BodyParser(loginDto); err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	email := loginDto.User.Email
	password := loginDto.User.Password

	userRepo := repository.NewUserRepo(database.GetDB())
	user, err := userRepo.GetByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return CreateErrorResponse(c, fiber.StatusUnauthorized, "Wrong email or password")
		} else {
			return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
		}
	}

	isValidPassword := IsValidPassword([]byte(user.Password), []byte(password))

	if isValidPassword != true {
		return CreateErrorResponse(c, fiber.StatusUnauthorized, "Wrong email or password")
	}

	token, err := GenerateAccessToken(user)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(fiber.Map{"user": fiber.Map{
		"email":    user.Email,
		"token":    token,
		"username": user.Username,
		"bio":      user.Bio,
		"image":    user.Image,
	}})
}

func GenerateHashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func IsValidPassword(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return false
	}

	return true
}

func GenerateAccessToken(user *model.User) (string, error) {

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
		},
	)
	s, err := t.SignedString([]byte(config.AppCfg().JWTSecretKey))
	if err != nil {
		return "", err
	}

	return s, err
}
