package routes

import (
	"visit-tracker/handlers"

	"github.com/gofiber/fiber/v3"
)

type AuthRoutes struct {
	handler *handlers.AuthHandler
}

func NewAuthRoutes(handler *handlers.AuthHandler) *AuthRoutes {
	return &AuthRoutes{handler: handler}
}

func (r *AuthRoutes) Register(app *fiber.App) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/register", r.handler.CreateUser)
	auth.Post("/verify", r.handler.VerifyEmail)
	auth.Post("/login", r.handler.Login)
}
