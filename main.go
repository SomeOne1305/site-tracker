package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"visit-tracker/config"
	"visit-tracker/database"
	"visit-tracker/handlers"
	"visit-tracker/mailer"
	"visit-tracker/repository"
	"visit-tracker/routes"
	"visit-tracker/services"
	"visit-tracker/storage"

	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.LoadConfig()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// CLIENTS
	redisClient := storage.NewRedisClient(cfg)
	mailerClient := mailer.NewMailer(cfg.MailerHost, cfg.MailerPort, cfg.MailerUsername, cfg.MailerPassword)

	authRepo := repository.NewAuthRepository(pool)
	authService := services.NewAuthService(authRepo, redisClient.Client, mailerClient)
	sessionRepo := repository.NewSessionRepository(pool)
	sessionService := services.NewSessionService(sessionRepo)
	authHandler := handlers.NewAuthHandler(authService, sessionService)
	authRoutes := routes.NewAuthRoutes(authHandler)

	userRepo := repository.NewUserRepository(pool)
	userHandler := handlers.NewUserHandler(userRepo)
	userRoutes := routes.NewUserRoutes(userHandler)

	projectRepo := repository.NewProjectRepository(pool)
	projectService := services.NewProjectRepository(*projectRepo)
	projectHandler := handlers.NewProjectHandler(*projectService)
	projectRoutes := routes.NewProjectRoutes(projectHandler)

	app := fiber.New()
	authRoutes.Register(app)
	userRoutes.Register(app, pool)
	projectRoutes.Register(app, pool)

	// Mount the UI with the default configuration under /swagger
	app.Get("/swagger/*", swaggo.HandlerDefault)
	// Customize the UI by passing a Config
	app.Get("/docs/*", swaggo.New(swaggo.Config{
		URL:               "http://example.com/doc.json",
		DeepLinking:       false,
		DocExpansion:      "none",
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	log.Printf("server listening on %s", cfg.ServerPort)
	if err := app.Listen(StdPort(cfg.ServerPort)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func StdPort(port string) string {
	if port == "" {
		return ":8080"
	}
	if port[0] != ':' {
		return ":" + port
	}
	return port
}
