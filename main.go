package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"petapp/internal/database"
	"petapp/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()
	defer database.DB.Close()

	app := fiber.New()

	routes.SetupUserRoutes(app)
	routes.SetupAuthRoutes(app)
	routes.SetupCategoryRoutes(app)
	routes.SetupProductRoutes(app)
	routes.SetupOrderRoutes(app)
	routes.SetupClinicRoutes(app)
	routes.SetupAppointmentRoutes(app)
	routes.SetupPetRoutes(app)
	routes.SetupAddressRoutes(app)
	routes.SetupClinicAuthRoutes(app)
	routes.SetupAdminAuthRoutes(app)
	routes.SetupAdminRoutes(app)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"
	}

	log.Fatal(app.Listen(port))
}
