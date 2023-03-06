package handlers

import (
	"com.ak.gooverlord/indexer"
	"com.ak.gooverlord/models"
	"github.com/gofiber/fiber/v2"
)

func Logs(c *fiber.Ctx) error {
	var entries []models.LogEntry
	if err := c.BodyParser(&entries); err != nil {
		return err
	}

	indexer.IndexEntries(entries)

	return c.JSON(fiber.Map{})
}
