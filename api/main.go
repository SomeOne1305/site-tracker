package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/valyala/fasthttp"
)

var app *fiber.App

func init() {
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

	trackingRepo := repository.NewTrackingRepository(pool)
	trackingService := services.NewTrackingService(trackingRepo)
	trackingHandler := handlers.NewTrackingHandler(trackingService)
	trackingRoutes := routes.NewTrackingRoutes(trackingHandler, redisClient.Client)

	app = fiber.New()
	authRoutes.Register(app)
	userRoutes.Register(app, pool)
	projectRoutes.Register(app, pool)
	trackingRoutes.Register(app, pool)

	// Serve swagger.json from docs folder
	app.Get("/docs/*", swaggo.New(swaggo.Config{
		URL: "/swagger/doc.json",
	}))
	// Serve swagger.json file
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Visit Tracker API is running")
	})
	app.Use(cors.New(cors.Config{
		AllowMethods: []string{
			"GET",
			"POST",
			"OPTIONS",
			"HEAD",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-API-Key",
		},
	}))
	app.Use(func(c fiber.Ctx) error { fmt.Println(c.Request().String()); return c.Next() })
}

// Handler is the Vercel serverless function entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	// Convert standard HTTP request to FastHTTP format for Fiber
	fastHTTPHandler := app.Handler()
	
	// Create FastHTTP request context
	ctx := &fasthttp.RequestCtx{}
	
	// Copy the request from net/http to fasthttp
	ctx.Request.Header.SetMethod(r.Method)
	ctx.Request.SetRequestURI(r.RequestURI)
	if r.RequestURI == "" {
		ctx.Request.SetRequestURI(r.URL.RequestURI())
	}
	
	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			ctx.Request.Header.Add(key, value)
		}
	}
	
	// Copy body if present
	if r.Body != nil {
		body, _ := io.ReadAll(r.Body)
		ctx.Request.SetBody(body)
	}
	
	// Handle the request with Fiber
	fastHTTPHandler(ctx)
	
	// Copy response headers
	ctx.Response.Header.VisitAll(func(key, value []byte) {
		w.Header().Set(string(key), string(value))
	})
	
	// Write status code and body
	w.WriteHeader(ctx.Response.StatusCode())
	w.Write(ctx.Response.Body())
}
