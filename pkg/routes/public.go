package routes

import (
	"github.com/fahrurben/realworld-go/app/controller"
	"github.com/gofiber/fiber/v2"
)

func Public(app *fiber.App) {
	app.Post("/api/users", controller.Register)
	app.Post("/api/users/login", controller.Login)
}
