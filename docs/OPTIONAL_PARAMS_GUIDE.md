# Руководство по использованию опциональных параметров в Queries

## Обзор изменений

Все параметры в методах `List`, `Paginate` и `Count` теперь являются **опциональными**. Это означает, что вы можете:
- Передавать `NULL` значения для параметров, которые не нужны
- Использовать только те фильтры, которые необходимы
- Не беспокоиться о передаче пустых значений

## Типы данных

### Опциональные типы из pgx/v5/pgtype

| SQL тип | Go тип | Описание |
|---------|---------|----------|
| `boolean` | `pgtype.Bool` | Опциональный boolean |
| `text` | `pgtype.Text` | Опциональный текст |
| `uuid[]` | `[]pgtype.UUID` | Массив UUID (nil = NULL) |
| `text[]` | `[]string` | Массив строк (nil = NULL) |

### Структуры параметров

#### ListAllRolesParams / ListAllPermissionsParams
```go
type ListAllRolesParams struct {
    ShowDeleted pgtype.Bool   `json:"show_deleted"` // Показывать удаленные записи
    Search      pgtype.Text   `json:"search"`       // Поисковая строка
    Values      []string      `json:"values"`       // Фильтр по значениям (values)
    Ids         []pgtype.UUID `json:"ids"`          // Фильтр по ID
    SortBy      interface{}   `json:"sort_by"`      // Поле для сортировки
    SortOrder   interface{}   `json:"sort_order"`   // Направление сортировки (ASC/DESC)
}
```

#### PaginateAllRolesParams / PaginateAllPermissionsParams
```go
type PaginateAllRolesParams struct {
    ShowDeleted pgtype.Bool   `json:"show_deleted"`
    Search      pgtype.Text   `json:"search"`
    Values      []string      `json:"values"`
    Ids         []pgtype.UUID `json:"ids"`
    SortBy      interface{}   `json:"sort_by"`
    SortOrder   interface{}   `json:"sort_order"`
    Limit       int32         `json:"limit"`        // Обязательный параметр
    Offset      int32         `json:"offset"`       // Обязательный параметр
}
```

#### CountAllRolesParams / CountAllPermissionsParams
```go
type CountAllRolesParams struct {
    ShowDeleted pgtype.Bool   `json:"show_deleted"`
    Search      pgtype.Text   `json:"search"`
    Values      []string      `json:"values"`
    Ids         []pgtype.UUID `json:"ids"`
}
```

## Примеры использования

### 1. Базовый запрос без фильтров

```go
// Получить все активные (не удаленные) роли
roles, err := queries.ListAllRoles(ctx, generated.ListAllRolesParams{
    ShowDeleted: pgtype.Bool{Valid: false},  // NULL - использовать значение по умолчанию
    Search:      pgtype.Text{Valid: false},  // NULL - без поиска
    Values:      nil,                         // NULL - без фильтра по values
    Ids:         nil,                         // NULL - без фильтра по IDs
    SortBy:      nil,                         // NULL - сортировка по умолчанию
    SortOrder:   nil,                         // NULL - порядок по умолчанию
})
```

### 2. Поиск по тексту

```go
// Найти роли, содержащие "admin" в любом из текстовых полей
roles, err := queries.ListAllRoles(ctx, generated.ListAllRolesParams{
    ShowDeleted: pgtype.Bool{Valid: false},
    Search:      pgtype.Text{String: "admin", Valid: true}, // Поиск по "admin"
    Values:      nil,
    Ids:         nil,
    SortBy:      "created_at",
    SortOrder:   "DESC",
})
```

### 3. Показать удаленные записи

```go
// Получить все роли, включая удаленные
roles, err := queries.ListAllRoles(ctx, generated.ListAllRolesParams{
    ShowDeleted: pgtype.Bool{Bool: true, Valid: true}, // Показать удаленные
    Search:      pgtype.Text{Valid: false},
    Values:      nil,
    Ids:         nil,
    SortBy:      nil,
    SortOrder:   nil,
})
```

### 4. Фильтр по значениям (values)

```go
// Получить только роли ADMIN и USER
roles, err := queries.ListAllRoles(ctx, generated.ListAllRolesParams{
    ShowDeleted: pgtype.Bool{Valid: false},
    Search:      pgtype.Text{Valid: false},
    Values:      []string{"ADMIN", "USER"}, // Фильтр по values
    Ids:         nil,
    SortBy:      "value",
    SortOrder:   "ASC",
})
```

### 5. Фильтр по ID

```go
// Получить роли с конкретными ID
id1 := pgtype.UUID{}
id1.Scan("123e4567-e89b-12d3-a456-426614174000")

id2 := pgtype.UUID{}
id2.Scan("123e4567-e89b-12d3-a456-426614174001")

roles, err := queries.ListAllRoles(ctx, generated.ListAllRolesParams{
    ShowDeleted: pgtype.Bool{Valid: false},
    Search:      pgtype.Text{Valid: false},
    Values:      nil,
    Ids:         []pgtype.UUID{id1, id2}, // Фильтр по IDs
    SortBy:      nil,
    SortOrder:   nil,
})
```

### 6. Пагинация с фильтрами

```go
// Получить первые 10 ролей со страницы 1 (offset 0)
roles, err := queries.PaginateAllRoles(ctx, generated.PaginateAllRolesParams{
    ShowDeleted: pgtype.Bool{Valid: false},
    Search:      pgtype.Text{String: "manager", Valid: true},
    Values:      []string{"MANAGER"},
    Ids:         nil,
    SortBy:      "created_at",
    SortOrder:   "DESC",
    Limit:       10,
    Offset:      0,
})
```

### 7. Подсчет записей

```go
// Посчитать все роли, соответствующие критериям
count, err := queries.CountAllRoles(ctx, generated.CountAllRolesParams{
    ShowDeleted: pgtype.Bool{Bool: true, Valid: true},
    Search:      pgtype.Text{String: "admin", Valid: true},
    Values:      nil,
    Ids:         nil,
})
```

## Helper функции

Для упрощения работы с опциональными параметрами можно создать вспомогательные функции:

```go
// Helper function для создания опционального Bool
func OptionalBool(value bool) pgtype.Bool {
    return pgtype.Bool{Bool: value, Valid: true}
}

// Helper function для создания опционального Text
func OptionalText(value string) pgtype.Text {
    return pgtype.Text{String: value, Valid: true}
}

// Helper function для создания NULL Bool
func NullBool() pgtype.Bool {
    return pgtype.Bool{Valid: false}
}

// Helper function для создания NULL Text
func NullText() pgtype.Text {
    return pgtype.Text{Valid: false}
}

// Использование:
roles, err := queries.ListAllRoles(ctx, generated.ListAllRolesParams{
    ShowDeleted: OptionalBool(true),
    Search:      OptionalText("admin"),
    Values:      []string{"ADMIN"},
    Ids:         nil,
    SortBy:      "created_at",
    SortOrder:   "DESC",
})
```

## Поведение по умолчанию

Когда параметр имеет значение `NULL` (или `Valid: false`), применяется следующее поведение:

| Параметр | Значение по умолчанию |
|----------|------------------------|
| `ShowDeleted` | `false` - скрывать удаленные записи |
| `Search` | Нет фильтрации по тексту |
| `Values` | Нет фильтрации по значениям |
| `Ids` | Нет фильтрации по ID |
| `SortBy` | `created_at` (сортировка по дате создания) |
| `SortOrder` | `DESC` (от новых к старым) |

## Миграция со старого API

### Было (с обязательными параметрами):
```go
roles, err := queries.ListAllRoles(ctx, generated.ListAllRolesParams{
    Column1: false,              // show_deleted
    Column2: "",                 // search
    Column3: []string{},         // values
    Column4: []pgtype.UUID{},    // ids
    Column5: "created_at",       // sort_by
    Column6: "DESC",             // sort_order
})
```

### Стало (с опциональными параметрами):
```go
roles, err := queries.ListAllRoles(ctx, generated.ListAllRolesParams{
    ShowDeleted: pgtype.Bool{Valid: false},    // NULL
    Search:      pgtype.Text{Valid: false},    // NULL
    Values:      nil,                           // NULL
    Ids:         nil,                           // NULL
    SortBy:      nil,                           // NULL
    SortOrder:   nil,                           // NULL
})
```

## Преимущества нового подхода

1. **Более читаемый код**: Именованные параметры вместо `Column1`, `Column2`, и т.д.
2. **Опциональность**: Не нужно передавать пустые значения для неиспользуемых фильтров
3. **Гибкость**: Легко комбинировать разные фильтры
4. **Типобезопасность**: Использование pgtype обеспечивает правильную работу с NULL значениями
5. **Самодокументируемость**: Понятно, какой параметр за что отвечает

## Технические детали

### SQL queries используют sqlc.narg() и sqlc.arg()

- `sqlc.narg('name')` - создает nullable именованный аргумент
- `sqlc.arg('name')` - создает обязательный именованный аргумент

Пример из SQL:
```sql
WHERE
    (CASE WHEN sqlc.narg('show_deleted')::boolean THEN TRUE ELSE r.deleted_at IS NULL END)
    AND (
        sqlc.narg('search')::text IS NULL OR
        r.title_ru ILIKE '%' || sqlc.narg('search') || '%'
    )
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset')
```

## Дополнительные примеры

Полные примеры использования можно найти в файле:
- `examples/optional_params_example.go`

## Заключение

Теперь все параметры в `List`, `Paginate`, и `Count` методах являются опциональными, что делает API более гибким и удобным в использовании. Вы можете передавать только те параметры, которые действительно нужны для вашего запроса.
