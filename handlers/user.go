package handlers

import (
	"fmt"
	"visit-tracker/models"
	"visit-tracker/repository"
	"visit-tracker/utils"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
	userID := c.Locals("userID")
	fmt.Println(c.Locals("userID"))
	user, err := h.repo.GetUser(c.Context(), userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User data is ready to use !", "user": user})
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {

	u := c.Locals("userID")
	fmt.Println(c.Locals("userID"))
	userID, _ := u.(string)
	fmt.Print("User id: ", userID)
	var req models.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		fmt.Println("Error binding request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if errs := utils.ValidateStruct(&req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_errors": errs})
	}

	user, err := h.repo.UpdateUser(c.Context(), req, userID)

	if err != nil {
		fmt.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to update user data"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}
