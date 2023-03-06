package main

import (
	"com.ak.gooverlord/handlers"
	_ "com.ak.gooverlord/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Post("/logs", handlers.Logs)
	app.Listen(":3000")
}
