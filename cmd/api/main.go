package main

import (
	"googleforms/handlers"
	"googleforms/internal/config"
	"googleforms/internal/database"
	"googleforms/internal/models"
	"googleforms/repositories"
	"googleforms/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Migrate tables
	db.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.Form{},
		&models.FormVersion{},
		&models.Submission{},
	)

	// Initialize repositories
	tenantRepo := repositories.NewTenantRepository(db)
	userRepo := repositories.NewUserRepository(db)
	formRepo := repositories.NewFormRepository(db)
	formVersionRepo := repositories.NewFormVersionRepository(db)
	submissionRepo := repositories.NewSubmissionRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, "your-jwt-secret")
	formService := services.NewFormService(formRepo, formVersionRepo, tenantRepo)
	submissionService := services.NewSubmissionService(submissionRepo, formRepo, formVersionRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Auth routes
	app.Post("/api/auth/login", authHandler.Login)
	app.Post("/api/auth/register", authHandler.Register)

	log.Println("   Server running on :8080")
	log.Println("   Available endpoints:")
	log.Println("   GET  /health")
	log.Println("   POST /api/auth/login")
	log.Println("   POST /api/auth/register")

	// error ..
	_ = formService
	_ = submissionService

	app.Listen(":8080")
}
