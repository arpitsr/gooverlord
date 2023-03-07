package handlers

import (
	"com.ak.gooverlord/models"
	"com.ak.gooverlord/query"
	"github.com/gofiber/fiber/v2"
)

func Query(c *fiber.Ctx) error {
	var searchQuery models.Query
	if err := c.BodyParser(&searchQuery); err != nil {
		return err
	}
	results := query.Search(searchQuery)
	return c.JSON(results)
}
