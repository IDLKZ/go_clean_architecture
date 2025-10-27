# Clean Architecture Fiber

A Clean Architecture implementation using Go Fiber framework with PostgreSQL and SQLC.

## Project Structure

```
clean_architecture_fiber/
├── app/                    # Application layer (handlers, middleware)
├── cmd/                    # Command-line tools
│   └── migrate/           # Database migration CLI
├── config/                 # Configuration management
├── core/                   # Core business logic
├── data/                   # Data layer
│   ├── db/
│   │   ├── generated/     # SQLC generated code
│   │   ├── migrate/       # Database migrations
│   │   ├── queries/       # SQL queries for SQLC
│   │   └── migrator.go    # Migration runner
│   └── entities/          # Entity definitions
├── domain/                 # Domain layer (interfaces, models)
└── shared/                 # Shared utilities
```

## Features

- ✅ Clean Architecture pattern
- ✅ Go Fiber web framework
- ✅ PostgreSQL database
- ✅ SQLC for type-safe SQL
- ✅ Database migrations with golang-migrate
- ✅ Soft delete support
- ✅ Multilingual support (Russian, English, Kazakh)
- ✅ Role-based access control (RBAC)
- ✅ Advanced filtering and pagination
- ✅ Dynamic sorting

## Prerequisites

- Go 1.24.9 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Quick Start

### 1. Clone the repository

```bash
git clone <repository-url>
cd clean_architecture_fiber
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Setup database

Create a PostgreSQL database:

```sql
CREATE DATABASE clean_architecture_db;
```

### 4. Configure environment

Set your database URL:

```bash
# Linux/Mac
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/clean_architecture_db?sslmode=disable"

# Windows (PowerShell)
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/clean_architecture_db?sslmode=disable"

# Windows (CMD)
set DATABASE_URL=postgres://postgres:postgres@localhost:5432/clean_architecture_db?sslmode=disable
```

Or edit the `DATABASE_URL` in Makefile.

### 5. Run migrations

```bash
# Using Makefile (recommended)
make schema-up

# Or directly
go run cmd/schema/main.go -database-url="$DATABASE_URL" -command=up
```

## Database Migrations

### Available Commands

```bash
# Show all available commands
make help

# Run all pending migrations
make schema-up

# Rollback last migration
make schema-down

# Run N migration steps (positive for up, negative for down)
make schema-steps STEPS=2
make schema-steps STEPS=-1

# Show current migration version
make schema-version

# Force migration version (use with caution!)
make schema-force VERSION=3

# Create new migration files
make schema-create NAME=add_users_table

# Setup database (run all migrations)
make db-setup

# Reset database (drop all + schema)
make db-reset
```

### Migration Files

Migrations are located in `data/db/migrate/`:

1. `000001_initial_db` - PostgreSQL extensions (pgcrypto)
2. `000002_create_roles_table` - Roles table
3. `000003_create_permissions_table` - Permissions table
4. `000004_create_role_permissions_table` - Role-Permission junction table

See [Migration Documentation](data/db/schema/README.md) for detailed information.

## Database Queries (SQLC)

### Generating Code

After modifying queries in `data/db/queries/`, regenerate Go code:

```bash
make sqlc-generate
```

### Available Queries

#### **Roles**
- `CreateOneRole` - Create single role
- `GetRoleById` - Get role by ID with permissions
- `GetRoleByValue` - Get role by value with permissions
- `UpdateRoleById` - Update role
- `DeleteRoleById` - Soft delete role
- `HardDeleteRoleById` - Permanently delete role
- `BulkCreateRoles` - Bulk insert roles
- `BulkUpdateRoles` - Bulk update roles
- `BulkDeleteRoleByIds` - Bulk soft delete
- `BulkHardDeleteRoleByIds` - Bulk hard delete
- `ListAllRoles` - List with filters and sorting
- `PaginateAllRoles` - Paginated list with filters
- `CountAllRoles` - Count with filters

#### **Permissions**
Similar operations as Roles, with joins to roles instead of permissions.

#### **Role-Permissions**
- `CreateOneRolePermission` - Create association
- `GetRolePermissionById` - Get with full role and permission data
- `DeleteRolePermissionById` - Delete association
- `BulkAssignPermissionsToRole` - Assign multiple permissions to role
- `BulkAssignRolesToPermission` - Assign multiple roles to permission
- `ListAllRolePermissions` - List with advanced filtering
- `PaginateAllRolePermissions` - Paginated list
- `CheckRoleHasPermission` - Check if role has permission
- And more...

### Query Features

**Filtering:**
- `show_deleted` - Show/hide soft-deleted records
- `search` - Full-text search across multiple fields
- `values` - Filter by array of values
- `ids` - Filter by array of IDs

**Sorting:**
- `order_by` - Field to sort by (created_at, updated_at, title_ru, value, etc.)
- `order_direction` - ASC or DESC

**Pagination:**
- `limit` - Page size
- `offset` - Starting position

### Example Usage

```go
import "clean_architecture_fiber/data/db/generated"

// List roles with search and sorting
roles, err := queries.ListAllRoles(ctx, generated.ListAllRolesParams{
    Column1: false,           // show_deleted
    Column2: "admin",         // search
    Column3: nil,             // values
    Column4: nil,             // ids
    Column5: "created_at",    // order_by
    Column6: "DESC",          // order_direction
})

// Paginate roles
roles, err := queries.PaginateAllRoles(ctx, generated.PaginateAllRolesParams{
    Column1: false,
    Column2: "",
    Column3: nil,
    Column4: nil,
    Column5: "title_ru",
    Column6: "ASC",
    Limit:   20,
    Offset:  0,
})
```

## Development

### Running the application

```bash
go run main.go
```

### Building the application

```bash
go build -o bin/app main.go
```

### Building migration tool

```bash
make build-schema
./bin/schema -database-url="..." -command=up
```

## Testing

```bash
go test ./...
```

## Database Schema

### Tables

**roles**
- Multilingual support (ru, en, kk)
- Soft delete support
- Unique value constraint
- Indexed on value

**permissions**
- Same structure as roles
- Multilingual support
- Soft delete support

**role_permissions**
- Junction table for many-to-many relationship
- Composite unique constraint on (role_id, permission_id)
- CASCADE delete on foreign keys

## Technologies

- **Web Framework**: [Fiber](https://github.com/gofiber/fiber)
- **Database**: PostgreSQL
- **SQL Generator**: [SQLC](https://sqlc.dev/)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Database Driver**: [pgx/v5](https://github.com/jackc/pgx)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Validation**: [validator/v10](https://github.com/go-playground/validator)

## License

[Your License Here]

## Contributing

[Contributing Guidelines Here]
