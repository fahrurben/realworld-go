package routes

import (
	"github.com/fahrurben/realworld-go/app/controller"
	"github.com/fahrurben/realworld-go/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func Protected(app *fiber.App) {
	api := app.Group("/api/users", middleware.JWTProtected())
	api.Get("/", controller.GetCurrentUser)
	api.Put("/", controller.UpdateUser)
}
