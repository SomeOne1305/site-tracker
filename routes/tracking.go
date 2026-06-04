package routes

import (
	"os"
	"time"
	"visit-tracker/handlers"
	"visit-tracker/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TrackingRoutes struct {
	handler     *handlers.TrackingHandler
	redisClient *redis.Client
}

func NewTrackingRoutes(handler *handlers.TrackingHandler, redisClient *redis.Client) *TrackingRoutes {
	return &TrackingRoutes{
		handler:     handler,
		redisClient: redisClient,
	}
}

func (r *TrackingRoutes) Register(app *fiber.App, pool *pgxpool.Pool) {
	// ===== PUBLIC ENDPOINTS (NO AUTH) =====
	app.Get("/embed/:projectID", r.handler.GetEmbedScript)
	app.Get("/t/:projectID.js", r.handler.GetScriptFile)
	app.Post("/api/track", r.rateLimitTracking(), r.handler.TrackPageView)

	// ===== PROTECTED ENDPOINTS (AUTH REQUIRED) =====
	protected := app.Group("/api/v1/projects/:projectID", middlewares.AuthRequired(pool))
	protected.Get("/dashboard", r.handler.GetDashboard)
	protected.Get("/chart/daily", r.handler.GetChartDaily)
	protected.Get("/chart/monthly", r.handler.GetChartMonthly)
	protected.Get("/endpoints", r.handler.GetAPIEndpoints)
	// ===== ADMIN ENDPOINTS (API KEY OR AUTH) =====
	admin := app.Group("/api/v1/admin")
	admin.Post("/aggregate", r.handler.TriggerAggregation)
}

// rateLimitTracking - Rate limiting middleware
func (r *TrackingRoutes) rateLimitTracking() fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID := c.FormValue("project_id")
		if projectID == "" {
			projectID = c.Query("project_id")
		}
		ip := c.IP()

		key := "track:" + projectID + ":" + ip
		count, _ := r.redisClient.Incr(c.Context(), key).Result()

		if count == 1 {
			r.redisClient.Expire(c.Context(), key, 60*time.Second)
		}

		if count > 200 {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "rate limited"})
		}

		return c.Next()
	}
}

// apiKeyAuth - Optional API key or bearer token auth for cron jobs
func (r *TrackingRoutes) apiKeyAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey != "" && apiKey == os.Getenv("CRON_API_KEY") {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
}
