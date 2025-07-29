
package routes

import (
	"petapp/internal/handlers"
	"petapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupClinicAuthRoutes(app *fiber.App) {
	clinicAuthGroup := app.Group("/api/clinic/auth")
	clinicAuthGroup.Post("/login", handlers.ClinicLogin)
	clinicAuthGroup.Get("/profile", middleware.ClinicAuthRequired, handlers.GetClinicProfile)
	clinicAuthGroup.Put("/profile", middleware.ClinicAuthRequired, handlers.UpdateClinicProfile)

	clinicAppointmentsGroup := app.Group("/api/clinic/appointments")
	clinicAppointmentsGroup.Use(middleware.ClinicAuthRequired)
	clinicAppointmentsGroup.Get("/", handlers.GetClinicAppointments)
	clinicAppointmentsGroup.Put("/:id/confirm", handlers.ConfirmAppointment)
	clinicAppointmentsGroup.Put("/:id/cancel", handlers.CancelClinicAppointment)
	clinicAppointmentsGroup.Get("/calendar", handlers.GetClinicAppointmentsCalendar)

	clinicScheduleGroup := app.Group("/api/clinic/schedule")
	clinicScheduleGroup.Use(middleware.ClinicAuthRequired)
	clinicScheduleGroup.Get("/", handlers.GetClinicSchedule)
	clinicScheduleGroup.Put("/", handlers.UpdateClinicSchedule)
}
