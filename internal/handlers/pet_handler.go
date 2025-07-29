package handlers

import (
	"petapp/internal/db"
	"petapp/internal/database"
	"petapp/internal/models"

	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CreatePetRequest struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	BirthDate string `json:"birth_date"`
}

type UpdatePetRequest struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	BirthDate string `json:"birth_date"`
}

func CreatePet(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	req := new(CreatePetRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid birth date format",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	pet, err := queries.CreatePet(c.Context(), db.CreatePetParams{
		Uuid:      sql.NullString{String: uuid.New().String(), Valid: true},
		UserID:    sql.NullInt64{Int64: userID, Valid: true},
		Name:      sql.NullString{String: req.Name, Valid: true},
		Type:      sql.NullString{String: req.Type, Valid: true},
		BirthDate: sql.NullTime{Time: birthDate, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not add pet",
			"error":   err.Error(),
		})
	}

	petResponse := models.PetResponse{
		ID:        pet.ID,
		UUID:      uuid.MustParse(pet.Uuid.String),
		UserID:    pet.UserID.Int64,
		Name:      pet.Name.String,
		Type:      pet.Type.String,
		BirthDate: pet.BirthDate.Time,
		CreatedAt: pet.CreatedAt.Time,
	}

	return c.Status(fiber.StatusCreated).JSON(petResponse)
}

func GetPets(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	queries := db.New(database.DB.DB())
	pets, err := queries.ListPetsByUserID(c.Context(), sql.NullInt64{Int64: userID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get user pets",
			"error":   err.Error(),
		})
	}

	petResponses := make([]models.PetResponse, len(pets))
	for i, pet := range pets {
		petResponses[i] = models.PetResponse{
			ID:        pet.ID,
			UUID:      uuid.MustParse(pet.Uuid.String),
			UserID:    pet.UserID.Int64,
			Name:      pet.Name.String,
			Type:      pet.Type.String,
			BirthDate: pet.BirthDate.Time,
			CreatedAt: pet.CreatedAt.Time,
		}
	}

	return c.JSON(petResponses)
}

func GetPetDetails(c *fiber.Ctx) error {
	petID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid pet ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	pet, err := queries.GetPetByID(c.Context(), petID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get pet details",
			"error":   err.Error(),
		})
	}

	petResponse := models.PetResponse{
		ID:        pet.ID,
		UUID:      uuid.MustParse(pet.Uuid.String),
		UserID:    pet.UserID.Int64,
		Name:      pet.Name.String,
		Type:      pet.Type.String,
		BirthDate: pet.BirthDate.Time,
		CreatedAt: pet.CreatedAt.Time,
	}

	return c.JSON(petResponse)
}

func UpdatePet(c *fiber.Ctx) error {
	petID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid pet ID",
			"error":   err.Error(),
		})
	}

	req := new(UpdatePetRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	params := db.UpdatePetParams{
		ID: petID,
	}

	if req.Name != "" {
		params.Name = sql.NullString{String: req.Name, Valid: true}
	}
	if req.Type != "" {
		params.Type = sql.NullString{String: req.Type, Valid: true}
	}
	if req.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid birth date format",
				"error":   err.Error(),
			})
		}
		params.BirthDate = sql.NullTime{Time: birthDate, Valid: true}
	}

	updatedPet, err := queries.UpdatePet(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update pet",
			"error":   err.Error(),
		})
	}

	petResponse := models.PetResponse{
		ID:        updatedPet.ID,
		UUID:      uuid.MustParse(updatedPet.Uuid.String),
		UserID:    updatedPet.UserID.Int64,
		Name:      updatedPet.Name.String,
		Type:      updatedPet.Type.String,
		BirthDate: updatedPet.BirthDate.Time,
		CreatedAt: updatedPet.CreatedAt.Time,
	}

	return c.JSON(petResponse)
}

func DeletePet(c *fiber.Ctx) error {
	petID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid pet ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	err = queries.DeletePet(c.Context(), petID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete pet",
			"error":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}