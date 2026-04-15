package handlers

import (
	"fmt"
	"strings"
	"time"
	"visit-tracker/models"
	"visit-tracker/services"
	"visit-tracker/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service *services.AuthService
	session *services.SessionService
}

func NewAuthHandler(service *services.AuthService, session *services.SessionService) *AuthHandler {
	return &AuthHandler{service: service, session: session}
}

func (h *AuthHandler) CreateUser(c fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if errs := utils.ValidateStruct(&req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_errors": errs})
	}

	user, err := h.service.Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AuthHandler) VerifyEmail(c fiber.Ctx) error {
	var req models.VerifyEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		fmt.Println("Error binding request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if errs := utils.ValidateStruct(&req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_errors": errs})
	}

	user, err := h.service.VerifyEmail(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if errs := utils.ValidateStruct(&req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_errors": errs})
	}

	user, err := h.service.Login(c.Context(), req)
	if err != nil {
		fmt.Println("Login error:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	refresh_token, _ := utils.GenerateRefreshToken()
	user_agent := c.Get("User-Agent")
	ip_address := c.IP()

	session, err := h.session.CreateSession(c.Context(), user.ID, refresh_token, user_agent, ip_address, time.Now().Add(time.Hour*24*3))
	if err != nil {
		fmt.Println("Session creation error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create session"})
	}

	sessionID := session.ID.String()
	soft := strings.Join([]string{sessionID, user.ID.String()}, "/")
	refresh_token_hashed, err := utils.EncryptID(soft)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create session"})
	}
	access_token, _ := utils.GenerateJWT(jwt.MapClaims{"id": user.ID.String(), "session_id": sessionID})

	c.Cookie(&fiber.Cookie{Expires: time.Now().Add(time.Hour * 24 * 3), HTTPOnly: true, Name: "refresh_token", Value: refresh_token_hashed})

	c.Cookie(&fiber.Cookie{Expires: time.Now().Add(time.Minute * 15), HTTPOnly: true, Name: "access_token", Value: access_token})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "login successful", "user": user})
}
