package main

import (
	"authentication/config"
	"authentication/handlers"
	"authentication/middleware"
	"authentication/types"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	config.InitDB()

	// Auto migrate database schema
	if err := config.DB.AutoMigrate(&types.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()

	// Public routes
	app.Post("/api/auth/signup", authHandler.Signup)
	app.Post("/api/auth/login", authHandler.Login)
	app.Post("/api/auth/forgot-password", authHandler.ForgotPassword)
	app.Post("/api/auth/reset-password", authHandler.ResetPassword)

	// Protected routes
	protected := app.Group("/api", middleware.AuthMiddleware())
	protected.Post("/auth/change-password", authHandler.ChangePassword)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
