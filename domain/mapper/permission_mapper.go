package mapper

import (
	"clean_architecture_fiber/data/db/generated"
	"clean_architecture_fiber/domain/dto"
	"time"

	"github.com/gofiber/fiber/v2"
)

// PermissionRDTOFromPermissionSQLC преобразует generated.Permission (sqlc) в dto.PermissionRDTO
// Автоматически выбирает нужный язык на основе текущего запроса
// Если перевод для запрошенного языка отсутствует, используется русский (fallback)
func PermissionRDTOFromPermissionSQLC(ctx *fiber.Ctx, permissionSQLC generated.Permission) dto.PermissionRDTO {
	// Получаем локализованные title и description
	title := getLocalizedText(ctx, permissionSQLC.TitleRu, permissionSQLC.TitleEn, permissionSQLC.TitleKk)
	description := getLocalizedText(ctx, permissionSQLC.DescriptionRu, permissionSQLC.DescriptionEn, permissionSQLC.DescriptionKk)

	// Преобразуем UUID в строку
	permissionID := uuidToString(permissionSQLC.ID)

	// Обрабатываем timestamps с учетом их валидности
	var createdAt, updatedAt time.Time
	var deletedAt *time.Time

	if permissionSQLC.CreatedAt.Valid {
		createdAt = permissionSQLC.CreatedAt.Time
	}

	if permissionSQLC.UpdatedAt.Valid {
		updatedAt = permissionSQLC.UpdatedAt.Time
	}

	// DeletedAt может быть nil (если разрешение не удалено)
	if permissionSQLC.DeletedAt.Valid {
		t := permissionSQLC.DeletedAt.Time
		deletedAt = &t
	}

	return dto.PermissionRDTO{
		ID:          permissionID,
		Title:       title,
		Value:       permissionSQLC.Value,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		DeletedAt:   deletedAt,
	}
}
