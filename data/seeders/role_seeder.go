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

// SeedRole инициализирует базовые роли в базе данных (ADMIN и MODERATOR)
// Проверяет наличие ролей и создает их только если база данных пуста
func SeedRole(ctx context.Context, pool *pgxpool.Pool) error {
	q := generated.New(pool)

	// Подсчитываем количество ролей (включая удаленные)
	arg := generated.CountAllRolesParams{
		ShowDeleted: pgtype.Bool{Bool: true, Valid: true}, // Проверяем все роли, включая удаленные
		Search:      pgtype.Text{Valid: false},            // Без поискового фильтра
		Values:      nil,                                  // Без фильтра по values
		Ids:         nil,                                  // Без фильтра по IDs
	}

	count, err := q.CountAllRoles(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to count roles: %w", err)
	}

	// Если роли уже есть, пропускаем создание
	if count > 0 {
		log.Printf("Roles already seeded. Found %d roles in database", count)
		return nil
	}

	// Создаем слайс для bulk insert (инициализируем с нужной длиной)
	parameters := make([]generated.BulkCreateRolesParams, 2)

	// Генерируем UUID для роли администратора
	adminUUID := uuid.New()
	adminPgUUID := pgtype.UUID{}
	if err := adminPgUUID.Scan(adminUUID.String()); err != nil {
		return fmt.Errorf("failed to create admin UUID: %w", err)
	}

	// Роль администратора
	parameters[0] = generated.BulkCreateRolesParams{
		ID:            adminPgUUID,
		TitleRu:       "Администратор",
		TitleEn:       pgtype.Text{String: "Administrator", Valid: true},
		TitleKk:       pgtype.Text{String: "Әкімші", Valid: true}, // Правильный перевод на казахский
		DescriptionRu: "Глобальная управляющая роль с доступом ко всем возможностям системы",
		DescriptionEn: pgtype.Text{String: "Global administrative role_use_case with access to all system features", Valid: true},
		DescriptionKk: pgtype.Text{String: "Жүйенің барлық мүмкіндіктеріне қолжетімділігі бар жаһандық басқарушы рөл", Valid: true},
		Value:         db_constants.AdminRoleValueConstant,
	}

	// Генерируем UUID для роли модератора
	moderatorUUID := uuid.New()
	moderatorPgUUID := pgtype.UUID{}
	if err := moderatorPgUUID.Scan(moderatorUUID.String()); err != nil {
		return fmt.Errorf("failed to create moderator UUID: %w", err)
	}

	// Роль модератора
	parameters[1] = generated.BulkCreateRolesParams{
		ID:            moderatorPgUUID,
		TitleRu:       "Модератор",
		TitleEn:       pgtype.Text{String: "Moderator", Valid: true},
		TitleKk:       pgtype.Text{String: "Модератор", Valid: true},
		DescriptionRu: "Роль модератора с ограниченным набором разрешений",
		DescriptionEn: pgtype.Text{String: "Moderator role_use_case with limited set of permissions", Valid: true},
		DescriptionKk: pgtype.Text{String: "Шектеулі рұқсаттар жиынтығы бар модератор рөлі", Valid: true},
		Value:         db_constants.ModeratorRoleValueConstant,
	}

	// Создаем роли через bulk insert
	_, err = q.BulkCreateRoles(ctx, parameters)
	if err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	// Логируем успешное создание
	log.Printf("Successfully seeded %d roles", len(parameters))

	return nil
}
