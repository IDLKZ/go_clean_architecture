package seeders

import (
	"clean_architecture_fiber/data/db/generated"
	"clean_architecture_fiber/shared/db_constants"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SeedRolePermission создает связи между ролями и разрешениями
// Назначает все базовые разрешения (MANAGE, CREATE, READ, EDIT, DELETE) роли администратора
func SeedRolePermission(ctx context.Context, pool *pgxpool.Pool) error {
	q := generated.New(pool)

	// Helper функция для создания UUID
	createUUID := func() (pgtype.UUID, error) {
		id := uuid.New()
		pgUUID := pgtype.UUID{}
		if err := pgUUID.Scan(id.String()); err != nil {
			return pgtype.UUID{}, fmt.Errorf("failed to create UUID: %w", err)
		}
		return pgUUID, nil
	}

	// Helper функция для проверки и создания связи роль-разрешение
	// Возвращает true если связь была создана, false если уже существовала
	createRolePermissionIfNotExists := func(roleValue string, permissionValue string) (bool, error) {
		// Получаем роль по значению
		role, err := q.GetRoleByValue(ctx, roleValue)
		if err != nil {
			return false, fmt.Errorf("failed to get role_use_case '%s': %w", roleValue, err)
		}

		// Получаем разрешение по значению
		permission, err := q.GetPermissionByValue(ctx, permissionValue)
		if err != nil {
			return false, fmt.Errorf("failed to get permission '%s': %w", permissionValue, err)
		}

		// Проверяем, существует ли уже такая связь
		rolePermissionCount, err := q.CountAllRolePermissions(ctx, generated.CountAllRolePermissionsParams{
			RoleIds:          nil,                       // Не фильтруем по ID
			PermissionIds:    nil,                       // Не фильтруем по ID
			RoleValues:       []string{roleValue},       // Фильтруем по значению роли
			PermissionValues: []string{permissionValue}, // Фильтруем по значению разрешения
		})
		if err != nil {
			return false, fmt.Errorf("failed to count role_use_case-permission for '%s'-'%s': %w", roleValue, permissionValue, err)
		}

		// Если связь уже существует, пропускаем создание
		if rolePermissionCount > 0 {
			log.Printf("Role-Permission link already exists: %s - %s", role.TitleRu, permission.TitleRu)
			return false, nil
		}

		// Генерируем UUID для новой связи
		newUUID, err := createUUID()
		if err != nil {
			return false, fmt.Errorf("failed to create UUID for role_use_case-permission: %w", err)
		}

		// Создаем связь роль-разрешение
		rolePermission, err := q.CreateOneRolePermission(ctx, generated.CreateOneRolePermissionParams{
			ID:           newUUID,
			RoleID:       role.ID,
			PermissionID: permission.ID,
		})
		if err != nil {
			return false, fmt.Errorf("failed to create role_use_case-permission for '%s'-'%s': %w", roleValue, permissionValue, err)
		}

		log.Printf("Successfully created role_use_case-permission: %s - %s (UUID: %s)", role.TitleRu, permission.TitleRu, rolePermission.ID.Bytes)
		return true, nil
	}

	// Список всех базовых разрешений для назначения администратору
	allPermissions := []string{
		db_constants.ManagePermissionPrefixConstant, // Полное управление CRUD
		db_constants.CreatePermissionPrefixConstant, // Создание записей
		db_constants.ReadPermissionPrefixConstant,   // Чтение записей
		db_constants.EditPermissionPrefixConstant,   // Редактирование записей
		db_constants.DeletePermissionPrefixConstant, // Удаление записей
	}

	// Счетчик созданных связей
	createdCount := 0

	// Проходим по каждому разрешению и создаем связь с ролью ADMIN
	for _, permissionValue := range allPermissions {
		created, err := createRolePermissionIfNotExists(db_constants.AdminRoleValueConstant, permissionValue)
		if err != nil {
			return fmt.Errorf("failed to process role_use_case-permission for ADMIN-%s: %w", permissionValue, err)
		}
		if created {
			createdCount++
		}
	}

	// Логируем итоговый результат
	if createdCount > 0 {
		log.Printf("Successfully seeded %d role_use_case-permission links", createdCount)
	} else {
		log.Printf("All role_use_case-permission links already exist, no new links created")
	}

	return nil
}
