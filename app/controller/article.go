package controller

import (
	"database/sql"
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gosimple/slug"
)

func CreateArticle(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := int64(claims["user_id"].(float64))

	createArticleDto := &model.CreateArticleDto{}
	if err := c.BodyParser(createArticleDto); err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	articleDto := createArticleDto.Article
	article := &model.Article{
		Title:       articleDto.Title,
		Slug:        slug.Make(articleDto.Title),
		Description: articleDto.Description,
		Body:        articleDto.Body,
		AuthorID:    user_id,
	}
	tagList := articleDto.TagList

	tagRepo := repository.NewTagRepo(database.GetDB())

	for _, tag := range tagList {
		existingTag, err := tagRepo.Get(tag)
		if err != nil && err == sql.ErrNoRows {
			existingTag = nil
		}

		if existingTag == nil {
			tagRepo.Create(tag)
		}
	}

	articleRepo := repository.NewArticleRepo(database.GetDB())
	id, err := articleRepo.Create(article)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	created, err := articleRepo.Get(id)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	userRepo := repository.NewUserRepo(database.GetDB())
	author, err := userRepo.Get(user_id)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	created.Tags = tagList
	created.Author = &model.Author{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: false,
	}

	return c.JSON(created)
}
