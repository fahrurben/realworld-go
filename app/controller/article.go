package controller

import (
	"database/sql"
	"github.com/fahrurben/realworld-go/app/model"
	"github.com/fahrurben/realworld-go/app/repository"
	"github.com/fahrurben/realworld-go/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gookit/goutil/arrutil"
	"github.com/gosimple/slug"
	"strconv"
)

func CreateArticle(c *fiber.Ctx) error {
	var user_id int64 = 0
	if token := c.Locals("user"); token != nil {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		user_id = int64(claims["user_id"].(float64))
	} else {
		return CreateErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	createArticleDto := &model.SaveArticleDto{}
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

	articleRepo := repository.NewArticleRepo(database.GetDB())
	id, err := articleRepo.Create(article)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	tagRepo := repository.NewTagRepo(database.GetDB())

	for _, tag := range tagList {
		existingTag, err := tagRepo.Get(tag)
		if err != nil && err == sql.ErrNoRows {
			existingTag = nil
		}

		if existingTag == nil {
			_, err := tagRepo.Create(tag)
			if err != nil {
				return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
			}
		}
	}

	for _, tagName := range tagList {
		_, err := articleRepo.CreateArticleTag(id, tagName)
		if err != nil {
			return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
		}
	}

	created, err := articleRepo.Get(id, &user_id)

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

	return c.JSON(fiber.Map{"article": created})
}

func GetArticle(c *fiber.Ctx) error {
	var user_id int64 = 0
	if token := c.Locals("user"); token != nil {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		user_id = int64(claims["user_id"].(float64))
	}
	article_slug := c.Params("slug")

	articleRepo := repository.NewArticleRepo(database.GetDB())
	article, err := articleRepo.GetBySlug(article_slug, &user_id)

	if err != nil {
		if err == sql.ErrNoRows {
			return CreateErrorResponse(c, fiber.StatusNotFound, "Article not found")
		} else {
			return CreateErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
	}

	userRepo := repository.NewUserRepo(database.GetDB())
	author, err := userRepo.Get(article.AuthorID)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	tags, err := articleRepo.GetArticleTags(article.ID)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	isFollowing := false
	if user_id > 0 {
		isFollowing, _ = userRepo.IsFollowing(user_id, article.AuthorID)
	}

	article.Tags = tags
	article.Author = &model.Author{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: isFollowing,
	}

	return c.JSON(fiber.Map{"article": article})
}

func UpdateArticle(c *fiber.Ctx) error {
	var user_id int64 = 0
	if token := c.Locals("user"); token != nil {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		user_id = int64(claims["user_id"].(float64))
	} else {
		return CreateErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	article_slug := c.Params("slug")

	updateArticleDto := &model.SaveArticleDto{}
	if err := c.BodyParser(updateArticleDto); err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}
	articleDto := updateArticleDto.Article

	articleRepo := repository.NewArticleRepo(database.GetDB())
	article, err := articleRepo.GetBySlug(article_slug, &user_id)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	article.Title = articleDto.Title
	article.Slug = slug.Make(articleDto.Title)
	article.Description = articleDto.Description
	article.Body = articleDto.Body

	err = articleRepo.Update(article.ID, article)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	tagList := articleDto.TagList

	tagRepo := repository.NewTagRepo(database.GetDB())

	articleTags, err := articleRepo.GetArticleTags(article.ID)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	for _, tag := range tagList {
		existingTag, err := tagRepo.Get(tag)
		if err != nil && err == sql.ErrNoRows {
			existingTag = nil
		}

		if existingTag == nil {
			_, err := tagRepo.Create(tag)
			if err != nil {
				return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
			}
		}
	}

	for _, tagName := range tagList {
		if !arrutil.SliceHas(articleTags, tagName) {
			_, err := articleRepo.CreateArticleTag(article.ID, tagName)
			if err != nil {
				return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
			}
		}
	}

	for _, tagName := range articleTags {
		if !arrutil.SliceHas(tagList, tagName) {
			err := articleRepo.DeleteArticleTag(article.ID, tagName)
			if err != nil {
				return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
			}
		}
	}

	userRepo := repository.NewUserRepo(database.GetDB())
	author, err := userRepo.Get(user_id)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	isFollowing := false
	if user_id > 0 {
		isFollowing, _ = userRepo.IsFollowing(user_id, article.AuthorID)
	}

	article.Tags = tagList
	article.Author = &model.Author{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: isFollowing,
	}

	return c.JSON(fiber.Map{"article": article})
}

func DeleteArticle(c *fiber.Ctx) error {
	var user_id int64 = 0
	if token := c.Locals("user"); token != nil {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		user_id = int64(claims["user_id"].(float64))
	} else {
		return CreateErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	article_slug := c.Params("slug")

	articleRepo := repository.NewArticleRepo(database.GetDB())
	article, err := articleRepo.GetBySlug(article_slug, &user_id)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	if article.AuthorID != user_id {
		return CreateErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	articleTags, err := articleRepo.GetArticleTags(article.ID)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	for _, tagName := range articleTags {
		err := articleRepo.DeleteArticleTag(article.ID, tagName)
		if err != nil {
			return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
		}
	}

	err = articleRepo.Delete(article.ID)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(nil)
}

func FavoriteArticle(c *fiber.Ctx) error {
	var user_id int64 = 0
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id = int64(claims["user_id"].(float64))

	article_slug := c.Params("slug")

	articleRepo := repository.NewArticleRepo(database.GetDB())
	article, err := articleRepo.GetBySlug(article_slug, &user_id)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	if !article.Favorited {
		err = articleRepo.Favorited(article.ID, user_id)
		if err != nil {
			return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
		}
	}

	article, err = articleRepo.GetBySlug(article_slug, &user_id)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	userRepo := repository.NewUserRepo(database.GetDB())
	author, err := userRepo.Get(user_id)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	articleTags, err := articleRepo.GetArticleTags(article.ID)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	isFollowing := false
	if user_id > 0 {
		isFollowing, _ = userRepo.IsFollowing(user_id, article.AuthorID)
	}

	article.Tags = articleTags
	article.Author = &model.Author{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: isFollowing,
	}

	return c.JSON(fiber.Map{"article": article})
}

func UnfavoriteArticle(c *fiber.Ctx) error {
	var user_id int64 = 0
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id = int64(claims["user_id"].(float64))

	article_slug := c.Params("slug")

	articleRepo := repository.NewArticleRepo(database.GetDB())
	article, err := articleRepo.GetBySlug(article_slug, &user_id)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	if article.Favorited {
		err = articleRepo.Unfavorited(article.ID, user_id)
		if err != nil {
			return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
		}
	}

	article, err = articleRepo.GetBySlug(article_slug, &user_id)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	userRepo := repository.NewUserRepo(database.GetDB())
	author, err := userRepo.Get(user_id)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	articleTags, err := articleRepo.GetArticleTags(article.ID)
	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	isFollowing := false
	if user_id > 0 {
		isFollowing, _ = userRepo.IsFollowing(user_id, article.AuthorID)
	}

	article.Tags = articleTags
	article.Author = &model.Author{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: isFollowing,
	}

	return c.JSON(fiber.Map{"article": article})
}

func ListArticle(c *fiber.Ctx) error {
	var user_id int64 = 0
	if token := c.Locals("user"); token != nil {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		user_id = int64(claims["user_id"].(float64))
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	offset, err := strconv.Atoi(c.Query("offset", "0"))

	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")

	articleRepo := repository.NewArticleRepo(database.GetDB())
	articles, articleCount, err := articleRepo.List(int64(limit), int64(offset), &tag, &author, &favorited)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	userRepo := repository.NewUserRepo(database.GetDB())

	for _, article := range articles {
		author, err := userRepo.Get(user_id)

		articleTags, err := articleRepo.GetArticleTags(article.ID)
		if err != nil {
			return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
		}

		isFollowing := false
		if user_id > 0 {
			isFollowing, _ = userRepo.IsFollowing(user_id, article.AuthorID)
		}

		article.Tags = articleTags
		article.Author = &model.Author{
			Username:  author.Username,
			Bio:       author.Bio,
			Image:     author.Image,
			Following: isFollowing,
		}
	}

	return c.JSON(fiber.Map{"articles": articles, "articlesCount": articleCount})
}

func FeedArticle(c *fiber.Ctx) error {
	var user_id int64 = 0
	if token := c.Locals("user"); token != nil {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		user_id = int64(claims["user_id"].(float64))
	} else {
		return CreateErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	offset, err := strconv.Atoi(c.Query("offset", "0"))

	userRepo := repository.NewUserRepo(database.GetDB())
	authors, err := userRepo.GetFollowings(user_id)

	articleRepo := repository.NewArticleRepo(database.GetDB())
	articles, articleCount, err := articleRepo.Feed(int64(limit), int64(offset), authors)

	if err != nil {
		return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	for _, article := range articles {
		author, err := userRepo.Get(user_id)

		articleTags, err := articleRepo.GetArticleTags(article.ID)
		if err != nil {
			return CreateErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
		}

		isFollowing := false
		if user_id > 0 {
			isFollowing, _ = userRepo.IsFollowing(user_id, article.AuthorID)
		}

		article.Tags = articleTags
		article.Author = &model.Author{
			Username:  author.Username,
			Bio:       author.Bio,
			Image:     author.Image,
			Following: isFollowing,
		}
	}

	return c.JSON(fiber.Map{"articles": articles, "articlesCount": articleCount})
}
