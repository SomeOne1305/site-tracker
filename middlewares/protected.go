package middlewares

import (
	"fmt"
	"strings"
	"time"
	"visit-tracker/repository"
	"visit-tracker/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// middleware/auth.go
func AuthRequired(pool *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Try cookie first, then Authorization header
		if c.Get("X-API-Key") != "" {
			return APIAuthRequired(pool)(c)
		}
		token := c.Cookies("access_token")
		refresh_token := c.Cookies("refresh_token")
		if refresh_token != "" && token == "" {
			id, err := utils.DecryptID(refresh_token)
			if err == nil {
				sessionID := strings.Split(id, "/")[0]
				userID := strings.Split(id, "/")[1]
				session, err := repository.NewSessionRepository(pool).GetSessionById(&fiber.DefaultCtx{}, uuid.MustParse(sessionID))

				if err != nil || session == nil {
					fmt.Print("Middleware 1: ", err)

					c.Cookie(&fiber.Cookie{Name: "refresh_token", Value: "", Expires: time.Unix(0, 0), HTTPOnly: true, Path: "/", MaxAge: -1})

					c.Cookie(&fiber.Cookie{Name: "access_token", Value: "", Expires: time.Unix(0, 0), HTTPOnly: true, Path: "/", MaxAge: -1})

					return c.Status(401).JSON(fiber.Map{"error": "missing credentials 1"})
				}
				access_token, _ := utils.GenerateJWT(jwt.MapClaims{"id": userID, "session_id": sessionID})
				c.Cookie(&fiber.Cookie{Expires: time.Now().Add(time.Minute * 15), HTTPOnly: true, Name: "access_token", Value: access_token})

				c.Locals("sessionID", id)
				c.Locals("userID", userID)
				return c.Next()
			}
			c.Cookie(&fiber.Cookie{Name: "refresh_token", HTTPOnly: true, MaxAge: -1})
			return c.Status(401).JSON(fiber.Map{"error": "missing credentials"})
		}
		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing access token"})
		}

		// Validate JWT (assuming you have a validate function)
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		// Store user info in context for downstream handlers
		mapClaims := claims.(jwt.MapClaims)

		session, err := repository.NewSessionRepository(pool).GetSessionById(&fiber.DefaultCtx{}, uuid.MustParse(mapClaims["session_id"].(string)))

		if err != nil || session == nil {
			fmt.Print("Middleware 1: ", err)
			return c.Status(401).JSON(fiber.Map{"error": "missing credentials 1"})
		}

		c.Locals("userID", mapClaims["id"])
		c.Locals("sessionID", mapClaims["session_id"])

		return c.Next()
	}
}

func APIAuthRequired(pool *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("X-API-Key")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing access token"})
		}

		decryptedID, err := utils.DecryptID(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		uuid, err := uuid.Parse(decryptedID)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals("projectID", uuid.String())
		return c.Next()
	}
}
