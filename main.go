package main

import (
	"github.com/assignment-go/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // this is the default limit of 10MB
	})

	app.Use(cors.New(cors.Config{}))

	route.Setup(app)

	listenErr := app.Listen(":3000")
	if listenErr != nil {
		panic("failed to connect to port 3000")
	}
}
