package handlers

import (
	"petapp/internal/db"
	"petapp/internal/database"
	"petapp/internal/models"

	"github.com/gofiber/fiber/v2"
)

func GetCategories(c *fiber.Ctx) error {
	queries := db.New(database.DB.DB())

	categories, err := queries.ListCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get categories",
			"error":   err.Error(),
		})
	}

	categoryResponses := make([]models.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = models.CategoryResponse{
			ID:     category.ID,
			Name:    category.Name.String,
			ImgUrl:  category.ImgUrl.String,
		}
	}

	return c.JSON(categoryResponses)
}