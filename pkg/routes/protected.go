package routes

import (
	"github.com/fahrurben/realworld-go/app/controller"
	"github.com/fahrurben/realworld-go/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func Protected(app *fiber.App) {
	userApi := app.Group("/api/users", middleware.JWTProtected())
	userApi.Get("/", controller.GetCurrentUser)
	userApi.Put("/", controller.UpdateUser)

	followApi := app.Group("/api/profiles/:username/follow", middleware.JWTProtected())
	followApi.Post("/", controller.FollowUser)
	followApi.Delete("/", controller.UnfollowUser)

	articleApi := app.Group("/api/articles", middleware.JWTProtected())
	articleApi.Post("/", controller.CreateArticle)

	favoriteArticleApi := app.Group("/api/articles/:slug/favorite", middleware.JWTChecked())
	favoriteArticleApi.Post("/", controller.FavoriteArticle)
	favoriteArticleApi.Delete("/", controller.UnfavoriteArticle)
}
