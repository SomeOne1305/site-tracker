package routes

import (
	"visit-tracker/handlers"
	"visit-tracker/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRoutes struct {
	handler *handlers.ProjectHandler
}

func NewProjectRoutes(handler *handlers.ProjectHandler) *ProjectRoutes {
	return &ProjectRoutes{handler: handler}
}

func (r *ProjectRoutes) Register(app *fiber.App, pool *pgxpool.Pool) {
	projects := app.Group("/api/v1/project", middlewares.AuthRequired(pool))
	projects.Post("/create", r.handler.CreateProject)
	projects.Get("/all", r.handler.GetProjects)
}
