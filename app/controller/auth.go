package controller

import (
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/pkg/config"
	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

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
