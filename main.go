package main

import (
	"fmt"

	"github.com/assignment-go/route"
	"github.com/assignment-go/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config, loadError := util.LoadEnv()
	if loadError != nil {
		panic(loadError)
	}
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // this is the default limit of 10MB
	})

	app.Use(cors.New())

	route.Setup(app)

	listenErr := app.Listen(":8080")
	if listenErr != nil {
		panic(fmt.Sprintf("failed to connect to port %s", config.Port))
	}
}
