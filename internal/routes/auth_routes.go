
package routes

import (
	"petapp/internal/handlers"
	"petapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	authGroup := app.Group("/api/auth")

	authGroup.Post("/register", handlers.Register)
	authGroup.Post("/login", handlers.Login)
	authGroup.Post("/logout", handlers.Logout)

	authGroup.Get("/profile", middleware.AuthRequired, handlers.GetProfile)
	authGroup.Put("/profile", middleware.AuthRequired, handlers.UpdateProfile)
}
