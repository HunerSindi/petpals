package handlers

import (
	"petapp/internal/db"
	"petapp/internal/database"
	"petapp/internal/models"

	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CreateAddressRequest struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	IsDefault    bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	IsDefault    bool   `json:"is_default"`
}

func CreateAddress(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	req := new(CreateAddressRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	address, err := queries.CreateUserAddress(c.Context(), db.CreateUserAddressParams{
		UserID:       sql.NullInt64{Int64: userID, Valid: true},
		AddressLine1: sql.NullString{String: req.AddressLine1, Valid: true},
		AddressLine2: sql.NullString{String: req.AddressLine2, Valid: req.AddressLine2 != ""},
		City:         sql.NullString{String: req.City, Valid: true},
		State:        sql.NullString{String: req.State, Valid: true},
		PostalCode:   sql.NullString{String: req.PostalCode, Valid: true},
		IsDefault:    sql.NullBool{Bool: req.IsDefault, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not add address",
			"error":   err.Error(),
		})
	}

	addressResponse := models.UserAddressResponse{
		ID:           address.ID,
		UserID:       address.UserID.Int64,
		AddressLine1: address.AddressLine1.String,
		AddressLine2: address.AddressLine2.String,
		City:         address.City.String,
		State:        address.State.String,
		PostalCode:   address.PostalCode.String,
		IsDefault:    address.IsDefault.Bool,
	}

	return c.Status(fiber.StatusCreated).JSON(addressResponse)
}

func GetAddresses(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	queries := db.New(database.DB.DB())
	addresses, err := queries.ListUserAddressesByUserID(c.Context(), sql.NullInt64{Int64: userID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get user addresses",
			"error":   err.Error(),
		})
	}

	addressResponses := make([]models.UserAddressResponse, len(addresses))
	for i, address := range addresses {
		addressResponses[i] = models.UserAddressResponse{
			ID:           address.ID,
			UserID:       address.UserID.Int64,
			AddressLine1: address.AddressLine1.String,
			AddressLine2: address.AddressLine2.String,
			City:         address.City.String,
			State:        address.State.String,
			PostalCode:   address.PostalCode.String,
			IsDefault:    address.IsDefault.Bool,
		}
	}

	return c.JSON(addressResponses)
}

func UpdateAddress(c *fiber.Ctx) error {
	addressID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid address ID",
			"error":   err.Error(),
		})
	}

	req := new(UpdateAddressRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	params := db.UpdateUserAddressParams{
		ID: addressID,
	}

	if req.AddressLine1 != "" {
		params.AddressLine1 = sql.NullString{String: req.AddressLine1, Valid: true}
	}
	if req.AddressLine2 != "" {
		params.AddressLine2 = sql.NullString{String: req.AddressLine2, Valid: true}
	}
	if req.City != "" {
		params.City = sql.NullString{String: req.City, Valid: true}
	}
	if req.State != "" {
		params.State = sql.NullString{String: req.State, Valid: true}
	}
	if req.PostalCode != "" {
		params.PostalCode = sql.NullString{String: req.PostalCode, Valid: true}
	}
	params.IsDefault = sql.NullBool{Bool: req.IsDefault, Valid: true}

	updatedAddress, err := queries.UpdateUserAddress(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update address",
			"error":   err.Error(),
		})
	}

	addressResponse := models.UserAddressResponse{
		ID:           updatedAddress.ID,
		UserID:       updatedAddress.UserID.Int64,
		AddressLine1: updatedAddress.AddressLine1.String,
		AddressLine2: updatedAddress.AddressLine2.String,
		City:         updatedAddress.City.String,
		State:        updatedAddress.State.String,
		PostalCode:   updatedAddress.PostalCode.String,
		IsDefault:    updatedAddress.IsDefault.Bool,
	}

	return c.JSON(addressResponse)
}

func DeleteAddress(c *fiber.Ctx) error {
	addressID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid address ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	err = queries.DeleteUserAddress(c.Context(), addressID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete address",
			"error":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}