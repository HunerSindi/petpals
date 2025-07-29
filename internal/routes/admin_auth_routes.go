
package routes

import (
	"petapp/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminAuthRoutes(app *fiber.App) {
	adminAuthGroup := app.Group("/api/admin/auth")
	adminAuthGroup.Post("/login", handlers.AdminLogin)
}
