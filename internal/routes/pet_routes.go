
package routes

import (
	"petapp/internal/handlers"
	"petapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupPetRoutes(app *fiber.App) {
	petGroup := app.Group("/api/pets")
	petGroup.Use(middleware.AuthRequired)

	petGroup.Post("/", handlers.CreatePet)
	petGroup.Get("/", handlers.GetPets)
	petGroup.Get("/:id", handlers.GetPetDetails)
	petGroup.Put("/:id", handlers.UpdatePet)
	petGroup.Delete("/:id", handlers.DeletePet)
}
