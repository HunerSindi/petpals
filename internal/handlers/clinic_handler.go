
package handlers

import (
	"petapp/internal/db"
	"petapp/internal/database"
	"petapp/internal/models"

	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetClinics(c *fiber.Ctx) error {
	queries := db.New(database.DB.DB())

	clinics, err := queries.ListClinics(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get clinics",
			"error":   err.Error(),
		})
	}

	clinicResponses := make([]models.ClinicResponse, len(clinics))
	for i, clinic := range clinics {
		images, err := queries.ListClinicImagesByClinicID(c.Context(), sql.NullInt64{Int64: clinic.ID, Valid: true})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not get clinic images",
				"error":   err.Error(),
			})
		}

		imgUrls := make([]string, 0)
		for _, img := range images {
			if img.ImgUrl != "" {
				imgUrls = append(imgUrls, img.ImgUrl)
			}
		}

		clinicResponses[i] = models.ClinicResponse{
			ID:         clinic.ID,
			FirstName:  clinic.FirstName.String,
			LastName:   clinic.LastName.String,
			ClinicName: clinic.ClinicName.String,
			Email:      clinic.Email.String,
			OpenTime:   clinic.OpenTime.Time,
			CloseTime:  clinic.CloseTime.Time,
			Description: clinic.Description.String,
			CreatedAt:  clinic.CreatedAt.Time,
			Images:     imgUrls,
		}
	}

	return c.JSON(clinicResponses)
}

func GetClinicDetails(c *fiber.Ctx) error {
	clinicID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid clinic ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	clinic, err := queries.GetClinicByID(c.Context(), clinicID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get clinic details",
			"error":   err.Error(),
		})
	}

	images, err := queries.ListClinicImagesByClinicID(c.Context(), sql.NullInt64{Int64: clinic.ID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get clinic images",
			"error":   err.Error(),
		})
	}

	imgUrls := make([]string, 0)
	for _, img := range images {
		if img.ImgUrl != "" {
			imgUrls = append(imgUrls, img.ImgUrl)
		}
	}

	clinicResponse := models.ClinicResponse{
		ID:         clinic.ID,
		FirstName:  clinic.FirstName.String,
		LastName:   clinic.LastName.String,
		ClinicName: clinic.ClinicName.String,
		Email:      clinic.Email.String,
		OpenTime:   clinic.OpenTime.Time,
		CloseTime:  clinic.CloseTime.Time,
		Description: clinic.Description.String,
		CreatedAt:  clinic.CreatedAt.Time,
		Images:     imgUrls,
	}

	return c.JSON(clinicResponse)
}

func GetClinicAvailableSlots(c *fiber.Ctx) error {
	clinicID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid clinic ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	slots, err := queries.GetClinicAvailableSlots(c.Context(), sql.NullInt64{Int64: clinicID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get available slots",
			"error":   err.Error(),
		})
	}

	appointmentResponses := make([]models.AppointmentResponse, len(slots))
	for i, slot := range slots {
		appointmentResponses[i] = models.AppointmentResponse{
			ID:              slot.ID,
			UserID:          slot.UserID.Int64,
			ClinicID:        slot.ClinicID.Int64,
			PetID:           slot.PetID.Int64,
			AppointmentDate: slot.AppointmentDate.Time,
			AppointmentTime: slot.AppointmentTime.Time,
			Status:          models.AppointmentStatus(slot.Status.AppointmentStatus),
			CreatedAt:       slot.CreatedAt.Time,
		}
	}

	return c.JSON(appointmentResponses)
}
