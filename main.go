package main

import (
	"log"

	"com.ak.gooverlord/handlers"
	_ "com.ak.gooverlord/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()
	app.Use(recover.New())
	app.Post("/logs", handlers.Logs)
	app.Post("/query", handlers.Query)
	err := app.Listen(":3000")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
