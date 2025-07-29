package handlers

import (
	"petapp/internal/db"
	"petapp/internal/database"
	"petapp/internal/models"
	"petapp/internal/utils"

	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetProductsByCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid category ID",
			"error":   err.Error(),
		})
	}

	limit, err := strconv.ParseInt(c.Query("limit", "10"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid limit parameter",
			"error":   err.Error(),
		})
	}

	offset, err := strconv.ParseInt(c.Query("offset", "0"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid offset parameter",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	products, err := queries.ListProductsByCategory(c.Context(), db.ListProductsByCategoryParams{
		CategoryID: sql.NullInt64{Int64: categoryID, Valid: true},
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get products by category",
			"error":   err.Error(),
		})
	}

	productResponses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		images, err := queries.ListProductImages(c.Context(), sql.NullInt64{Int64: product.ID, Valid: true})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not get product images",
				"error":   err.Error(),
			})
		}

		imgUrls := make([]string, len(images))
		for j, img := range images {
			imgUrls[j] = img.ImgUrl.String
		}

		productResponses[i] = models.ProductResponse{
			ID:          product.ID,
			CategoryID:  product.CategoryID.Int64,
			Name:        product.Name.String,
			Description: product.Description.String,
			Price:       utils.ParseFloat(product.Price.String),
			CreatedAt:   product.CreatedAt.Time,
			Images:      imgUrls,
		}
	}

	return c.JSON(productResponses)
}

func GetProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	product, err := queries.GetProduct(c.Context(), productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get product details",
			"error":   err.Error(),
		})
	}

	images, err := queries.ListProductImages(c.Context(), sql.NullInt64{Int64: product.ID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get product images",
			"error":   err.Error(),
		})
	}

	imgUrls := make([]string, len(images))
	for i, img := range images {
		imgUrls[i] = img.ImgUrl.String
	}

	productResponse := models.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID.Int64,
		Name:        product.Name.String,
		Description: product.Description.String,
		Price:       utils.ParseFloat(product.Price.String),
		CreatedAt:   product.CreatedAt.Time,
		Images:      imgUrls,
	}

	return c.JSON(productResponse)
}

func GetProductImages(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	images, err := queries.ListProductImages(c.Context(), sql.NullInt64{Int64: productID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get product images",
			"error":   err.Error(),
		})
	}

	imageResponses := make([]models.ProductImageResponse, len(images))
	for i, image := range images {
		imageResponses[i] = models.ProductImageResponse{
			ID:        image.ID,
			ProductID: image.ProductID.Int64,
			ImgUrl:    image.ImgUrl.String,
			IsPrimary: image.IsPrimary.Bool,
		}
	}

	return c.JSON(imageResponses)
}