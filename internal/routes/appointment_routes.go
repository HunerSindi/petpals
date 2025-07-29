package routes

import (
	"petapp/internal/handlers"
	"petapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAppointmentRoutes(app *fiber.App) {
	appointmentGroup := app.Group("/api/appointments")
	appointmentGroup.Use(middleware.AuthRequired)

	appointmentGroup.Post("/", handlers.CreateAppointment)
	appointmentGroup.Get("/", handlers.GetAppointments)
	// appointmentGroup.Put("/:id", handlers.UpdateAppointment)
	appointmentGroup.Delete("/:id", handlers.DeleteAppointment)
}
