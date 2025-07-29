
package routes

import (
	"petapp/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupClinicRoutes(app *fiber.App) {
	clinicGroup := app.Group("/api/clinics")
	clinicGroup.Get("/", handlers.GetClinics)
	clinicGroup.Get("/:id", handlers.GetClinicDetails)
	clinicGroup.Get("/:id/available-slots", handlers.GetClinicAvailableSlots)
}
