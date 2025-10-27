# Database Migrations

This directory contains database migration files for the Clean Architecture Fiber project.

## Migration Files Structure

Migrations follow the naming convention: `{version}_{description}.{up|down}.sql`

Example:
- `000001_initial_db.up.sql` - Creates initial database setup
- `000001_initial_db.down.sql` - Reverts initial database setup

## Available Migrations

1. **000001_initial_db** - Initialize database with required extensions (pgcrypto)
2. **000002_create_roles_table** - Create roles table with multilingual support
3. **000003_create_permissions_table** - Create permissions table with multilingual support
4. **000004_create_role_permissions_table** - Create role-permission junction table

## Running Migrations

### Prerequisites

Install the required Go packages:
```bash
go get -u github.com/golang-schema/schema/v4
go get -u github.com/golang-schema/schema/v4/database/postgres
go get -u github.com/golang-schema/schema/v4/source/iofs
```

### Using Makefile (Recommended)

```bash
# Show available commands
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

# Drop all database tables (WARNING: destructive!)
make schema-drop

# Create new migration files
make schema-create NAME=add_users_table

# Setup database (run all migrations)
make db-setup

# Reset database (drop + schema)
make db-reset
```

### Using Go Command Directly

```bash
# Run all migrations
go run cmd/schema/main.go -database-url="postgres://user:pass@localhost:5432/dbname?sslmode=disable" -command=up

# Rollback last migration
go run cmd/schema/main.go -database-url="..." -command=down

# Run N steps
go run cmd/schema/main.go -database-url="..." -command=steps -steps=2

# Show version
go run cmd/schema/main.go -database-url="..." -command=version

# Force version
go run cmd/schema/main.go -database-url="..." -command=force -version=3
```

### Using Environment Variable

You can also set the `DATABASE_URL` environment variable:

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
make schema-up
```

Or:

```bash
# Windows (PowerShell)
$env:DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
make schema-up

# Windows (CMD)
set DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
make schema-up
```

## Creating New Migrations

### Using Makefile
```bash
make schema-create NAME=add_new_feature
```

This will create two files:
- `{timestamp}_add_new_feature.up.sql`
- `{timestamp}_add_new_feature.down.sql`

### Manual Creation

1. Create two files with the next version number:
   - `{version}_{description}.up.sql` - Contains the forward migration
   - `{version}_{description}.down.sql` - Contains the rollback migration

2. Write your SQL in the `.up.sql` file

3. Write the reverse SQL in the `.down.sql` file to undo the changes

## Best Practices

1. **Always create both up and down migrations** - This allows you to rollback changes if needed

2. **Keep migrations small and focused** - One migration should do one thing

3. **Test migrations thoroughly** - Test both up and down migrations before committing

4. **Never modify existing migrations** - Once a migration has been run in production, create a new migration to modify the schema

5. **Use transactions when possible** - Wrap your migrations in transactions for atomic operations

6. **Check migration status** - Always check `make migrate-version` before running migrations

7. **Backup before migrations** - Always backup your database before running migrations in production

## Migration States

- **Clean** - Migration completed successfully
- **Dirty** - Migration failed partway through and needs to be fixed

If a migration is in a dirty state:
1. Fix the migration file
2. Force the version back: `make migrate-force VERSION=X`
3. Re-run the migration: `make migrate-up`

## Example Migration

### Up Migration (`000005_add_users_table.up.sql`)
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    role_id UUID REFERENCES roles(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role_id ON users(role_id);
```

### Down Migration (`000005_add_users_table.down.sql`)
```sql
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

## Troubleshooting

### Migration fails with "dirty database"
```bash
# Check current version and status
make schema-version

# Force to the last known good version
make schema-force VERSION=X

# Try running migrations again
make schema-up
```

### Can't connect to database
- Verify your DATABASE_URL is correct
- Check that PostgreSQL is running
- Verify network connectivity
- Check firewall settings

### Migration already applied
This is normal - the migration system tracks which migrations have been applied and will skip them.

## Database Schema

After running all migrations, your database will have:

- **roles** table - User roles with multilingual support
- **permissions** table - Permission definitions with multilingual support
- **role_permissions** table - Many-to-many relationship between roles and permissions
- **pgcrypto** extension - For UUID generation and cryptographic functions
