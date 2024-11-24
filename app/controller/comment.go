package controller

import (
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func CreateComment(c *fiber.Ctx) error {
	var user_id int64 = 0
	if token := c.Locals("user"); token != nil {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		user_id = int64(claims["user_id"].(float64))
	} else {
		return CreateErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	article_slug := c.Params("slug")

	createCommentDto := &model.SaveCommentDTO{}
	if err := c.BodyParser(createCommentDto); err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	commentDto := createCommentDto.Comment

	articleRepo := repository.NewArticleRepo(database.GetDB())
	article, err := articleRepo.GetBySlug(article_slug, &user_id)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	commentRepo := repository.NewCommentRepository(database.GetDB())
	id, err := commentRepo.Create(article.ID, user_id, commentDto.Body)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	comment, err := commentRepo.Get(id)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	userRepo := repository.NewUserRepo(database.GetDB())
	author, err := userRepo.Get(comment.AuthorID)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	isFollowing := false
	if user_id > 0 {
		isFollowing, _ = userRepo.IsFollowing(user_id, article.AuthorID)
	}

	comment.Author = &model.Author{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: isFollowing,
	}

	return c.JSON(fiber.Map{"comment": comment})
}
