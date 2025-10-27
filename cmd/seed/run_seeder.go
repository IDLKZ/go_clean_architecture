package seed

import (
	"clean_architecture_fiber/data/seeders"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunSeeders выполняет все сидеры для инициализации базы данных
// Порядок выполнения важен:
// 1. SeedRole - создает базовые роли (ADMIN, MODERATOR)
// 2. SeedPermission - создает базовые разрешения (MANAGE, CREATE, READ, EDIT, DELETE)
// 3. SeedRolePermission - связывает роли с разрешениями (ADMIN получает все разрешения)
//
// При ошибке любого сидера приложение завершается с fatal error
func RunSeeders(ctx context.Context, pool *pgxpool.Pool) {
	// Шаг 1: Создание ролей
	// Создает роли ADMIN и MODERATOR если они еще не существуют
	if err := seeders.SeedRole(ctx, pool); err != nil {
		log.Fatalf("❌ SeedRole failed: %v", err)
	}

	// Шаг 2: Создание разрешений
	// Создает базовые разрешения: MANAGE, CREATE, READ, EDIT, DELETE
	if err := seeders.SeedPermission(ctx, pool); err != nil {
		log.Fatalf("❌ SeedPermission failed: %v", err)
	}

	// Шаг 3: Связывание ролей с разрешениями
	// Назначает все разрешения роли ADMIN
	if err := seeders.SeedRolePermission(ctx, pool); err != nil {
		log.Fatalf("❌ SeedRolePermission failed: %v", err)
	}

	// Все сидеры выполнены успешно
	log.Println("🌱 All seeders executed successfully!")
}
