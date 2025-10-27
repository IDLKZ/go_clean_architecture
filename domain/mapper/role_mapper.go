package mapper

import (
	"clean_architecture_fiber/data/db/generated"
	"clean_architecture_fiber/domain/dto"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RoleRDTOFromRoleSQLC преобразует generated.Role (sqlc) в dto.RoleRDTO
// Автоматически выбирает нужный язык на основе текущего запроса
// Если перевод для запрошенного языка отсутствует, используется русский (fallback)
func RoleRDTOFromRoleSQLC(ctx *fiber.Ctx, roleSQLC *generated.GetRoleByValueRow) *dto.RoleRDTO {
	if roleSQLC != nil {
		// Получаем локализованные title и description
		title := getLocalizedText(ctx, roleSQLC.TitleRu, roleSQLC.TitleEn, roleSQLC.TitleKk)
		description := getLocalizedText(ctx, roleSQLC.DescriptionRu, roleSQLC.DescriptionEn, roleSQLC.DescriptionKk)

		// Преобразуем UUID в строку
		roleID := uuidToString(roleSQLC.ID)

		// Обрабатываем timestamps с учетом их валидности
		var createdAt, updatedAt, deletedAt time.Time

		if roleSQLC.CreatedAt.Valid {
			createdAt = roleSQLC.CreatedAt.Time
		}

		if roleSQLC.UpdatedAt.Valid {
			updatedAt = roleSQLC.UpdatedAt.Time
		}

		// DeletedAt может быть nil (если роль не удалена)
		if roleSQLC.DeletedAt.Valid {
			deletedAt = roleSQLC.DeletedAt.Time
		}

		return &dto.RoleRDTO{
			ID:          roleID,
			Title:       title,
			Value:       roleSQLC.Value,
			Description: description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			DeletedAt:   deletedAt,
		}
	}
	return nil

}
