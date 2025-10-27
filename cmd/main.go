package main

import (
	"clean_architecture_fiber/cmd/seed"
	"clean_architecture_fiber/cmd/server"
	"clean_architecture_fiber/config"
	i18nPkg "clean_architecture_fiber/pkg/i18n"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// main является точкой входа в приложение
// Инициализирует конфигурацию, подключение к БД, i18n, запускает сидеры и Fiber сервер
func main() {
	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализируем систему локализации (i18n)
	log.Println("🌍 Initializing i18n...")
	if err := i18nPkg.Init(); err != nil {
		log.Fatalf("❌ Failed to initialize i18n: %v", err)
	}
	log.Printf("✅ i18n initialized (supported languages: %v)", i18nPkg.SupportedLanguages)

	// Загружаем конфигурацию приложения
	log.Println("📋 Loading application configuration...")
	cfg := config.LoadAppConfig()

	// Получаем DSN для подключения к базе данных
	dsn := cfg.GetDatabaseURL()
	log.Printf("🔌 Connecting to database: %s", cfg.Database.Name)

	// Создаем пул подключений к PostgreSQL
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Проверяем подключение к базе данных
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}
	log.Println("✅ Database connection established")

	// Запускаем сидеры для инициализации базовых данных
	log.Println("🌱 Running database seeders...")
	seed.RunSeeders(ctx, pool)

	// Настраиваем и запускаем Fiber сервер
	log.Printf("🚀 Starting %s server on port %d...", cfg.App.Name, cfg.App.Port)
	server.SettleFiberApp(cfg, ctx, pool)

	// Ожидание сигнала завершения (Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server gracefully...")
}
