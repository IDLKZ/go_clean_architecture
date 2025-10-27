package role_use_case

import (
	"clean_architecture_fiber/domain/dto"
	"clean_architecture_fiber/domain/mapper"
	"clean_architecture_fiber/domain/repositories"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type GetRoleByValueInput struct {
	Value string
}

type GetRoleByValueUseCase struct {
	Repo repositories.RoleRepository
}

func NewGetRoleByValueUseCase(repo repositories.RoleRepository) *GetRoleByValueUseCase {
	return &GetRoleByValueUseCase{Repo: repo}
}

// --- Реализация UseCase интерфейса ---

func (u *GetRoleByValueUseCase) Validate(fiberCtx *fiber.Ctx, ctx context.Context, input GetRoleByValueInput) error {
	if input.Value == "" {
		return errors.New("value cannot be empty")
	}
	return nil
}

func (u *GetRoleByValueUseCase) Execute(fiberCtx *fiber.Ctx, ctx context.Context, input GetRoleByValueInput) (*dto.RoleRDTO, error) {
	roleSQLC, err := u.Repo.GetByValue(ctx, input.Value)
	if err != nil {
		return nil, err
	}
	if roleSQLC == nil {
		return nil, fmt.Errorf("role not found")
	}
	return mapper.RoleRDTOFromRoleSQLC(fiberCtx, roleSQLC), nil
}

func (u *GetRoleByValueUseCase) Transform(fiberCtx *fiber.Ctx, ctx context.Context, result dto.RoleRDTO) (any, error) {
	return result, nil
}
