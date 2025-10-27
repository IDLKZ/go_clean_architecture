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

// SeedPermission инициализирует базовые разрешения в базе данных
// Создает 5 базовых разрешений: MANAGE, CREATE, READ, EDIT, DELETE
// Проверяет наличие разрешений и создает их только если база данных пуста
func SeedPermission(ctx context.Context, pool *pgxpool.Pool) error {
	q := generated.New(pool)

	// Подсчитываем количество разрешений (включая удаленные)
	countArgs := generated.CountAllPermissionsParams{
		ShowDeleted: pgtype.Bool{Bool: true, Valid: true}, // Проверяем все разрешения, включая удаленные
		Search:      pgtype.Text{Valid: false},            // Без поискового фильтра
		Values:      nil,                                   // Без фильтра по values
		Ids:         nil,                                   // Без фильтра по IDs
	}

	count, err := q.CountAllPermissions(ctx, countArgs)
	if err != nil {
		return fmt.Errorf("failed to count permissions: %w", err)
	}

	// Если разрешения уже есть, пропускаем создание
	if count > 0 {
		log.Printf("Permissions already seeded. Found %d permissions in database", count)
		return nil
	}

	// Создаем слайс для bulk insert (инициализируем с нужной длиной)
	parameters := make([]generated.BulkCreatePermissionsParams, 5)

	// Helper функция для создания UUID
	createUUID := func() (pgtype.UUID, error) {
		id := uuid.New()
		pgUUID := pgtype.UUID{}
		if err := pgUUID.Scan(id.String()); err != nil {
			return pgtype.UUID{}, fmt.Errorf("failed to create UUID: %w", err)
		}
		return pgUUID, nil
	}

	// 1. Разрешение MANAGE (Управление всем CRUD)
	manageUUID, err := createUUID()
	if err != nil {
		return err
	}
	parameters[0] = generated.BulkCreatePermissionsParams{
		ID:            manageUUID,
		TitleRu:       "Полное управление CRUD",
		TitleEn:       pgtype.Text{String: "Manage All CRUD", Valid: true},
		TitleKk:       pgtype.Text{String: "CRUD-ті толық басқару", Valid: true},
		DescriptionRu: "Глобальное управление всеми операциями: создание, чтение, редактирование и удаление",
		DescriptionEn: pgtype.Text{String: "Global management of all operations: create, read, update and delete", Valid: true},
		DescriptionKk: pgtype.Text{String: "Барлық операцияларды жаһандық басқару: жасау, оқу, өңдеу және жою", Valid: true},
		Value:         db_constants.ManagePermissionPrefixConstant,
	}

	// 2. Разрешение CREATE (Создание)
	createPermUUID, err := createUUID()
	if err != nil {
		return err
	}
	parameters[1] = generated.BulkCreatePermissionsParams{
		ID:            createPermUUID,
		TitleRu:       "Создание записей",
		TitleEn:       pgtype.Text{String: "Create All", Valid: true},
		TitleKk:       pgtype.Text{String: "Жазбаларды жасау", Valid: true},
		DescriptionRu: "Глобальное разрешение на создание всех типов записей",
		DescriptionEn: pgtype.Text{String: "Global permission to create all types of records", Valid: true},
		DescriptionKk: pgtype.Text{String: "Барлық жазба түрлерін жасауға жаһандық рұқсат", Valid: true},
		Value:         db_constants.CreatePermissionPrefixConstant,
	}

	// 3. Разрешение READ (Чтение)
	readUUID, err := createUUID()
	if err != nil {
		return err
	}
	parameters[2] = generated.BulkCreatePermissionsParams{
		ID:            readUUID,
		TitleRu:       "Чтение записей",
		TitleEn:       pgtype.Text{String: "Read All", Valid: true},
		TitleKk:       pgtype.Text{String: "Жазбаларды оқу", Valid: true},
		DescriptionRu: "Глобальное разрешение на чтение всех типов записей",
		DescriptionEn: pgtype.Text{String: "Global permission to read all types of records", Valid: true},
		DescriptionKk: pgtype.Text{String: "Барлық жазба түрлерін оқуға жаһандық рұқсат", Valid: true},
		Value:         db_constants.ReadPermissionPrefixConstant,
	}

	// 4. Разрешение EDIT (Редактирование)
	editUUID, err := createUUID()
	if err != nil {
		return err
	}
	parameters[3] = generated.BulkCreatePermissionsParams{
		ID:            editUUID,
		TitleRu:       "Редактирование записей",
		TitleEn:       pgtype.Text{String: "Edit All", Valid: true},
		TitleKk:       pgtype.Text{String: "Жазбаларды өңдеу", Valid: true},
		DescriptionRu: "Глобальное разрешение на редактирование всех типов записей",
		DescriptionEn: pgtype.Text{String: "Global permission to edit all types of records", Valid: true},
		DescriptionKk: pgtype.Text{String: "Барлық жазба түрлерін өңдеуге жаһандық рұқсат", Valid: true},
		Value:         db_constants.EditPermissionPrefixConstant,
	}

	// 5. Разрешение DELETE (Удаление)
	deleteUUID, err := createUUID()
	if err != nil {
		return err
	}
	parameters[4] = generated.BulkCreatePermissionsParams{
		ID:            deleteUUID,
		TitleRu:       "Удаление записей",
		TitleEn:       pgtype.Text{String: "Delete All", Valid: true},
		TitleKk:       pgtype.Text{String: "Жазбаларды жою", Valid: true},
		DescriptionRu: "Глобальное разрешение на удаление всех типов записей",
		DescriptionEn: pgtype.Text{String: "Global permission to delete all types of records", Valid: true},
		DescriptionKk: pgtype.Text{String: "Барлық жазба түрлерін жоюға жаһандық рұқсат", Valid: true},
		Value:         db_constants.DeletePermissionPrefixConstant,
	}

	// Создаем разрешения через bulk insert
	_, err = q.BulkCreatePermissions(ctx, parameters)
	if err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	// Логируем успешное создание
	log.Printf("Successfully seeded %d permissions", len(parameters))

	return nil
}
