package dto

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// PermissionDTO используется для операций создания/обновления разрешений
type PermissionDTO struct {
	ID          pgtype.UUID `json:"id,omitempty"`
	TitleRu     string      `json:"title_ru"`
	TitleEn     string      `json:"title_en"`
	TitleKk     string      `json:"title_kk"`
	Value       string      `json:"value"`
	Description string      `json:"description"`
	IsActive    bool        `json:"is_active"`
}

// PermissionRDTO используется для чтения (Read) разрешений с автоматической локализацией
// Title и Description автоматически выбираются на основе языка запроса
type PermissionRDTO struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Value       string     `json:"value"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
