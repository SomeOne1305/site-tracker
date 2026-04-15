package handlers

import (
	"fmt"
	"strings"
	"visit-tracker/models"
	"visit-tracker/services"
	"visit-tracker/utils"

	"github.com/gofiber/fiber/v3"
)

type ProjectHandler struct {
	service services.ProjectService
}

func NewProjectHandler(service services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) CreateProject(c fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)
	var req models.CreateProjectRequest
	if err := c.Bind().Body(&req); err != nil {
		fmt.Println("Error binding request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if errs := utils.ValidateStruct(&req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_errors": errs})
	}

	project, err := h.service.CreateProject(c.Context(), userID, req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": strings.Join([]string{"failed due to", err.Error()}, "")})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Project data is ready to use !", "user": project})
}

func (h *ProjectHandler) GetProjects(c fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)
	projects, err := h.service.GetProjects(c.Context(), userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": strings.Join([]string{"failed due to", err.Error()}, "")})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Project data is ready to use !", "data": projects})
}
