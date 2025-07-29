
package routes

import (
	"petapp/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/api/users")
	userGroup.Get("/", handlers.GetUsers)
}
