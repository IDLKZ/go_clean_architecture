package dto

import (
	"time"
)

// RoleRDTO используется для чтения (Read) ролей с автоматической локализацией
// Title и Description автоматически выбираются на основе языка запроса
type RoleRDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
}
