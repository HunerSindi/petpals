
package handlers

import (
	"petapp/internal/db"
	"petapp/internal/database"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	queries := db.New(database.DB.DB())

	users, err := queries.ListUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get users",
			"error":   err.Error(),
		})
	}

	return c.JSON(users)
}
