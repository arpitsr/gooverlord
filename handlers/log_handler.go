package handlers

import (
	"log"

	"com.ak.gooverlord/indexer"
	"com.ak.gooverlord/models"
	"github.com/gofiber/fiber/v2"
)

func Logs(c *fiber.Ctx) error {
	var entries []models.LogEntry
	if err := c.BodyParser(&entries); err != nil {
		log.Printf("%s", err)
		return err
	}

	indexer.GetInstance().IndexEntries(entries)

	return c.JSON(fiber.Map{})
}
