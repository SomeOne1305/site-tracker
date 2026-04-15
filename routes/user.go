package routes

import (
	"visit-tracker/handlers"
	"visit-tracker/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRoutes struct {
	handler *handlers.UserHandler
}

func NewUserRoutes(handler *handlers.UserHandler) *UserRoutes {
	return &UserRoutes{handler: handler}
}

func (r *UserRoutes) Register(app *fiber.App, pool *pgxpool.Pool) {
	userGroup := app.Group("/api/v1/user", middlewares.AuthRequired(pool))
	userGroup.Get("/me", r.handler.GetUser)
	userGroup.Put("update-me", r.handler.UpdateUser)
}
