package handlers

import (
	"os"
	"time"

	"petapp/internal/db"
	"petapp/internal/database"
	"petapp/internal/models"

	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ClinicLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateClinicProfileRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ClinicName string `json:"clinic_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	OpenTime   string `json:"open_time"`
	CloseTime  string `json:"close_time"`
	Description string `json:"description"`
}

func ClinicLogin(c *fiber.Ctx) error {
	req := new(ClinicLoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	clinic, err := queries.GetClinicByEmail(c.Context(), sql.NullString{String: req.Email, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(clinic.Password.String), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	expiresAt := time.Now().Add(time.Hour * 24).Unix()
	claims := jwt.MapClaims{
		"clinic_id": clinic.ID,
		"exp":       expiresAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate token",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{"token": signedToken})
}

func GetClinicProfile(c *fiber.Ctx) error {
	clinic := c.Locals("clinic").(*jwt.Token)
	claims := clinic.Claims.(jwt.MapClaims)
	clinicID := int64(claims["clinic_id"].(float64))

	queries := db.New(database.DB.DB())
	profile, err := queries.GetClinicByID(c.Context(), clinicID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get clinic profile",
			"error":   err.Error(),
		})
	}

	clinicResponse := models.ClinicResponse{
		ID:         profile.ID,
		FirstName:  profile.FirstName.String,
		LastName:   profile.LastName.String,
		ClinicName: profile.ClinicName.String,
		Email:      profile.Email.String,
		OpenTime:   profile.OpenTime.Time,
		CloseTime:  profile.CloseTime.Time,
		Description: profile.Description.String,
		CreatedAt:  profile.CreatedAt.Time,
	}

	return c.JSON(clinicResponse)
}

func UpdateClinicProfile(c *fiber.Ctx) error {
	clinic := c.Locals("clinic").(*jwt.Token)
	claims := clinic.Claims.(jwt.MapClaims)
	clinicID := int64(claims["clinic_id"].(float64))

	req := new(UpdateClinicProfileRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	params := db.UpdateClinicParams{
		ID: clinicID,
	}

	if req.FirstName != "" {
		params.FirstName = sql.NullString{String: req.FirstName, Valid: true}
	}
	if req.LastName != "" {
		params.LastName = sql.NullString{String: req.LastName, Valid: true}
	}
	if req.ClinicName != "" {
		params.ClinicName = sql.NullString{String: req.ClinicName, Valid: true}
	}
	if req.Email != "" {
		params.Email = sql.NullString{String: req.Email, Valid: true}
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not hash password",
				"error":   err.Error(),
			})
		}
		params.Password = sql.NullString{String: string(hashedPassword), Valid: true}
	}
	if req.OpenTime != "" {
		openTime, err := time.Parse("15:04", req.OpenTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid open time format",
				"error":   err.Error(),
			})
		}
		params.OpenTime = sql.NullTime{Time: openTime, Valid: true}
	}
	if req.CloseTime != "" {
		closeTime, err := time.Parse("15:04", req.CloseTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid close time format",
				"error":   err.Error(),
			})
		}
		params.CloseTime = sql.NullTime{Time: closeTime, Valid: true}
	}
	if req.Description != "" {
		params.Description = sql.NullString{String: req.Description, Valid: true}
	}

	updatedClinic, err := queries.UpdateClinic(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update clinic profile",
			"error":   err.Error(),
		})
	}

	clinicResponse := models.ClinicResponse{
		ID:         updatedClinic.ID,
		FirstName:  updatedClinic.FirstName.String,
		LastName:   updatedClinic.LastName.String,
		ClinicName: updatedClinic.ClinicName.String,
		Email:      updatedClinic.Email.String,
		OpenTime:   updatedClinic.OpenTime.Time,
		CloseTime:  updatedClinic.CloseTime.Time,
		Description: updatedClinic.Description.String,
		CreatedAt:  updatedClinic.CreatedAt.Time,
	}

	return c.JSON(clinicResponse)
}

func GetClinicAppointments(c *fiber.Ctx) error {
	clinic := c.Locals("clinic").(*jwt.Token)
	claims := clinic.Claims.(jwt.MapClaims)
	clinicID := int64(claims["clinic_id"].(float64))

	queries := db.New(database.DB.DB())
	appointments, err := queries.ListAppointmentsByClinicID(c.Context(), sql.NullInt64{Int64: clinicID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get clinic appointments",
			"error":   err.Error(),
		})
	}

	appointmentResponses := make([]models.AppointmentResponse, len(appointments))
	for i, appointment := range appointments {
		appointmentResponses[i] = models.AppointmentResponse{
			ID:              appointment.ID,
			UserID:          appointment.UserID.Int64,
			ClinicID:        appointment.ClinicID.Int64,
			PetID:           appointment.PetID.Int64,
			AppointmentDate: appointment.AppointmentDate.Time,
			AppointmentTime: appointment.AppointmentTime.Time,
			Status:          models.AppointmentStatus(appointment.Status.AppointmentStatus),
			CreatedAt:       appointment.CreatedAt.Time,
		}
	}

	return c.JSON(appointmentResponses)
}

func ConfirmAppointment(c *fiber.Ctx) error {
	appointmentID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid appointment ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	appointment, err := queries.UpdateAppointment(c.Context(), db.UpdateAppointmentParams{
		ID:     appointmentID,
		Status: db.NullAppointmentStatus{AppointmentStatus: db.AppointmentStatusConfirmed, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not confirm appointment",
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

	return c.JSON(appointmentResponse)
}

func CancelClinicAppointment(c *fiber.Ctx) error {
	appointmentID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid appointment ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	appointment, err := queries.UpdateAppointment(c.Context(), db.UpdateAppointmentParams{
		ID:     appointmentID,
		Status: db.NullAppointmentStatus{AppointmentStatus: db.AppointmentStatusCancelled, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not cancel appointment",
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

	return c.JSON(appointmentResponse)
}

func GetClinicAppointmentsCalendar(c *fiber.Ctx) error {
	clinic := c.Locals("clinic").(*jwt.Token)
	claims := clinic.Claims.(jwt.MapClaims)
	clinicID := int64(claims["clinic_id"].(float64))

	queries := db.New(database.DB.DB())
	// This query needs to be more sophisticated to return a calendar view.
	// For now, it will return all appointments for the clinic.
	appointments, err := queries.ListAppointmentsByClinicID(c.Context(), sql.NullInt64{Int64: clinicID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get clinic appointments for calendar view",
			"error":   err.Error(),
		})
	}

	appointmentResponses := make([]models.AppointmentResponse, len(appointments))
	for i, appointment := range appointments {
		appointmentResponses[i] = models.AppointmentResponse{
			ID:              appointment.ID,
			UserID:          appointment.UserID.Int64,
			ClinicID:        appointment.ClinicID.Int64,
			PetID:           appointment.PetID.Int64,
			AppointmentDate: appointment.AppointmentDate.Time,
			AppointmentTime: appointment.AppointmentTime.Time,
			Status:          models.AppointmentStatus(appointment.Status.AppointmentStatus),
			CreatedAt:       appointment.CreatedAt.Time,
		}
	}

	return c.JSON(appointmentResponses)
}

func GetClinicSchedule(c *fiber.Ctx) error {
	clinic := c.Locals("clinic").(*jwt.Token)
	claims := clinic.Claims.(jwt.MapClaims)
	clinicID := int64(claims["clinic_id"].(float64))

	queries := db.New(database.DB.DB())
	clinicDetails, err := queries.GetClinicByID(c.Context(), clinicID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get clinic schedule",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"open_time":  clinicDetails.OpenTime.Time.Format("15:04"),
		"close_time": clinicDetails.CloseTime.Time.Format("15:04"),
	})
}

func UpdateClinicSchedule(c *fiber.Ctx) error {
	clinic := c.Locals("clinic").(*jwt.Token)
	claims := clinic.Claims.(jwt.MapClaims)
	clinicID := int64(claims["clinic_id"].(float64))

	req := new(UpdateClinicProfileRequest) // Reusing for simplicity, consider dedicated struct
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	params := db.UpdateClinicScheduleParams{
		ID: clinicID,
	}

	if req.OpenTime != "" {
		openTime, err := time.Parse("15:04", req.OpenTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid open time format",
				"error":   err.Error(),
			})
		}
		params.OpenTime = sql.NullTime{Time: openTime, Valid: true}
	}
	if req.CloseTime != "" {
		closeTime, err := time.Parse("15:04", req.CloseTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid close time format",
				"error":   err.Error(),
			})
		}
		params.CloseTime = sql.NullTime{Time: closeTime, Valid: true}
	}

	updatedClinic, err := queries.UpdateClinicSchedule(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update clinic schedule",
			"error":   err.Error(),
		})
	}

	clinicResponse := models.ClinicResponse{
		ID:         updatedClinic.ID,
		FirstName:  updatedClinic.FirstName.String,
		LastName:   updatedClinic.LastName.String,
		ClinicName: updatedClinic.ClinicName.String,
		Email:      updatedClinic.Email.String,
		OpenTime:   updatedClinic.OpenTime.Time,
		CloseTime:  updatedClinic.CloseTime.Time,
		Description: updatedClinic.Description.String,
		CreatedAt:  updatedClinic.CreatedAt.Time,
	}

	return c.JSON(clinicResponse)
}