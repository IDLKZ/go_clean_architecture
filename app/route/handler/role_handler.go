package handler

import (
	"clean_architecture_fiber/domain/use_case/role_use_case"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	GetRoleByValueUC *role_use_case.GetRoleByValueUseCase
}

func NewRoleHandler(getUC *role_use_case.GetRoleByValueUseCase) *RoleHandler {
	return &RoleHandler{GetRoleByValueUC: getUC}
}

// GET /api/v1/roles/:value
func (h *RoleHandler) GetByValue(c *fiber.Ctx) error {
	value := c.Params("value")
	if value == "" {
		return fiber.NewError(http.StatusBadRequest, "role value is required")
	}

	result, err := h.GetRoleByValueUC.Execute(c, context.Background(), role_use_case.GetRoleByValueInput{Value: value})
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if result == nil {
		return fiber.NewError(http.StatusNotFound, "role not found")
	}

	return c.Status(http.StatusOK).JSON(result)
}
