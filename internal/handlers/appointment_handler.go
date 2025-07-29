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
)

type CreateAppointmentRequest struct {
	ClinicID       int64  `json:"clinic_id"`
	PetID          int64  `json:"pet_id"`
	AppointmentDate string `json:"appointment_date"`
	AppointmentTime string `json:"appointment_time"`
}

type UpdateAppointmentRequest struct {
	ClinicID       int64  `json:"clinic_id"`
	PetID          int64  `json:"pet_id"`
	AppointmentDate string `json:"appointment_date"`
	AppointmentTime string `json:"appointment_time"`
	Status          string `json:"status"`
}

func CreateAppointment(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	req := new(CreateAppointmentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid appointment date format",
			"error":   err.Error(),
		})
	}

	appointmentTime, err := time.Parse("15:04", req.AppointmentTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid appointment time format",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	appointment, err := queries.CreateAppointment(c.Context(), db.CreateAppointmentParams{
		UserID:          sql.NullInt64{Int64: userID, Valid: true},
		ClinicID:        sql.NullInt64{Int64: req.ClinicID, Valid: true},
		PetID:           sql.NullInt64{Int64: req.PetID, Valid: true},
		AppointmentDate: sql.NullTime{Time: appointmentDate, Valid: true},
		AppointmentTime: sql.NullTime{Time: appointmentTime, Valid: true},
		Status:          db.NullAppointmentStatus{AppointmentStatus: db.AppointmentStatusPending, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not book appointment",
			"error":   err.Error(),
		})
	}

	appointmentResponse := models.AppointmentResponse{
		ID:              appointment.ID,
		UserID:          appointment.UserID.Int64,
		ClinicID:        appointment.ClinicID.Int64,
		PetID:           appointment.PetID.Int64,
		AppointmentDate: appointment.AppointmentDate.Time,
		AppointmentTime: appointment.AppointmentTime.Time,
		Status:          models.AppointmentStatus(appointment.Status.AppointmentStatus),
		CreatedAt:       appointment.CreatedAt.Time,
	}

	return c.Status(fiber.StatusCreated).JSON(appointmentResponse)
}

func GetAppointments(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	queries := db.New(database.DB.DB())
	appointments, err := queries.ListAppointmentsDetailsByUserID(c.Context(), sql.NullInt64{Int64: userID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get user appointments",
			"error":   err.Error(),
		})
	}

	appointmentResponses := make([]models.AppointmentDetailsResponse, len(appointments))
	for i, appointment := range appointments {
		appointmentResponses[i] = models.AppointmentDetailsResponse{
			ID:              appointment.ID,
			AppointmentDate: appointment.AppointmentDate.Time,
			AppointmentTime: appointment.AppointmentTime.Time,
			Status:          models.AppointmentStatus(appointment.Status.AppointmentStatus),
			CreatedAt:       appointment.CreatedAt.Time,
			User: models.UserResponse{
				ID:        appointment.UserID,
				FirstName: appointment.FirstName.String,
				LastName:  appointment.LastName.String,
				Phone:     appointment.Phone.String,
				Email:     appointment.Email.String,
			},
			Clinic: models.ClinicResponse{
				ID:         appointment.ClinicID,
				ClinicName: appointment.ClinicName.String,
				Email:      appointment.ClinicEmail.String,
				OpenTime:   appointment.OpenTime.Time,
				CloseTime:  appointment.CloseTime.Time,
				Description: appointment.ClinicDescription.String,
			},
			Pet: models.PetResponse{
				ID:        appointment.PetID,
				Name:      appointment.PetName.String,
				Type:      appointment.PetType.String,
				BirthDate: appointment.PetBirthDate.Time,
			},
		}
	}

	return c.JSON(appointmentResponses)
}

func UpdateAppointment(c *fiber.Ctx) error {
	appointmentID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid appointment ID",
			"error":   err.Error(),
		})
	}

	req := new(UpdateAppointmentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	params := db.UpdateAppointmentParams{
		ID: appointmentID,
	}

	if req.ClinicID != 0 {
		params.ClinicID = sql.NullInt64{Int64: req.ClinicID, Valid: true}
	}
	if req.PetID != 0 {
		params.PetID = sql.NullInt64{Int64: req.PetID, Valid: true}
	}
	if req.AppointmentDate != "" {
		appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid appointment date format",
				"error":   err.Error(),
			})
		}
		params.AppointmentDate = sql.NullTime{Time: appointmentDate, Valid: true}
	}
	if req.AppointmentTime != "" {
		appointmentTime, err := time.Parse("15:04", req.AppointmentTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid appointment time format",
				"error":   err.Error(),
			})
		}
		params.AppointmentTime = sql.NullTime{Time: appointmentTime, Valid: true}
	}
	if req.Status != "" {
		params.Status = db.NullAppointmentStatus{AppointmentStatus: db.AppointmentStatus(req.Status), Valid: true}
	}

	updatedAppointment, err := queries.UpdateAppointment(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update appointment",
			"error":   err.Error(),
		})
	}

	appointmentResponse := models.AppointmentResponse{
		ID:              updatedAppointment.ID,
		UserID:          updatedAppointment.UserID.Int64,
		ClinicID:        updatedAppointment.ClinicID.Int64,
		PetID:           updatedAppointment.PetID.Int64,
		AppointmentDate: updatedAppointment.AppointmentDate.Time,
		AppointmentTime: updatedAppointment.AppointmentTime.Time,
		Status:          models.AppointmentStatus(updatedAppointment.Status.AppointmentStatus),
		CreatedAt:       updatedAppointment.CreatedAt.Time,
	}

	return c.JSON(appointmentResponse)
}

func DeleteAppointment(c *fiber.Ctx) error {
	appointmentID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid appointment ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	err = queries.DeleteAppointment(c.Context(), appointmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete appointment",
			"error":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}