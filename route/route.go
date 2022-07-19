package route

import (
	"github.com/assignment-go/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/upload", controller.Upload)
}
