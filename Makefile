# ⚙️ Переменные
DB_URL := $(shell go run cmd/tools/print_dsn.go)
MIGRATE_PATH := data/db/schema

# 🚀 Применить миграции
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

# 🌀 Откатить миграции
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

# 💥 Полный откат
migrate-reset:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down

# 🧩 Создать новую миграцию
migrate-create:
	@if [ -z "$(name)" ]; then echo "⚠️  Укажите имя миграции: make migrate-create name=create_users"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

# 🔍 Проверить версию
migrate-version:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" version
