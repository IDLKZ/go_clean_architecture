package main

import (
	"clean_architecture_fiber/cmd/seed"
	"clean_architecture_fiber/config"
	"clean_architecture_fiber/core/dependecy_injection"
	i18nPkg "clean_architecture_fiber/pkg/i18n"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

// main является точкой входа в приложение
// Инициализирует конфигурацию, подключение к БД, i18n, запускает сидеры и Fiber сервер
func main() {
	// Инициализируем систему локализации (i18n)
	log.Println("🌍 Initializing i18n...")
	if err := i18nPkg.Init(); err != nil {
		log.Fatalf("❌ Failed to initialize i18n: %v", err)
	}
	log.Printf("✅ i18n initialized (supported languages: %v)", i18nPkg.SupportedLanguages)

	// Создаем Fx приложение с DI контейнером
	app := fx.New(
		// Предоставляем конфигурацию
		fx.Provide(func() *config.Config {
			log.Println("📋 Loading application configuration...")
			return config.LoadAppConfig()
		}),
		// Подключаем основной модуль приложения
		dependecy_injection.AppModule,
		// Запускаем сидеры после инициализации БД
		fx.Invoke(func(lc fx.Lifecycle, pool *pgxpool.Pool) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Println("🌱 Running database seeders...")
					seed.RunSeeders(ctx, pool)
					return nil
				},
			})
		}),
	)

	// Запускаем приложение и ждем сигнала завершения
	app.Run()
}
