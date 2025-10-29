package dependecy_injection

import (
	"clean_architecture_fiber/app/route"
	"clean_architecture_fiber/config"
	"clean_architecture_fiber/data/db/generated"
	"context"
	"fmt"
	"log"

	i18nPkg "clean_architecture_fiber/pkg/i18n"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

// NewFiberApp создает и настраивает экземпляр Fiber приложения
func NewFiberApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:               cfg.Fiber.Prefork,
		CaseSensitive:         cfg.Fiber.CaseSensitive,
		StrictRouting:         cfg.Fiber.StrictRouting,
		ServerHeader:          cfg.App.Name,
		AppName:               cfg.App.Name,
		Concurrency:           cfg.Fiber.Concurrency,
		ReadTimeout:           cfg.Fiber.ReadTimeout,
		WriteTimeout:          cfg.Fiber.WriteTimeout,
		IdleTimeout:           cfg.Fiber.IdleTimeout,
		EnablePrintRoutes:     cfg.Fiber.EnablePrintRoutes,
		EnableIPValidation:    cfg.Fiber.EnableIPValidation,
		Immutable:             cfg.Fiber.Immutable,
		ProxyHeader:           cfg.Fiber.ProxyHeader,
		DisableStartupMessage: cfg.Fiber.DisableStartupMessage,
	})

	// Устанавливаем глобальные middleware
	setupMiddleware(app)

	return app
}

// setupMiddleware настраивает глобальные middleware для приложения
func setupMiddleware(app *fiber.App) {
	// Recover middleware - восстановление после паники
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Logger middleware - логирование HTTP запросов
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// RequestID middleware - добавление уникального ID к каждому запросу
	app.Use(requestid.New())

	// I18n middleware - определение языка запроса
	app.Use(i18nPkg.Middleware())

	// CORS middleware - настройка Cross-Origin Resource Sharing
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Accept-Language",
		AllowCredentials: false,
	}))

	// Compress middleware - сжатие ответов
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
}

// NewPgPool создает пул подключений к PostgreSQL
func NewPgPool(lc fx.Lifecycle, cfg *config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()

	dsn := cfg.GetDatabaseURL()
	log.Printf("🔌 Connecting to database: %s", cfg.Database.Name)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established")

	// Регистрируем хук для закрытия пула при остановке приложения
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("🔌 Closing database connection...")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

// StartFiberServer запускает Fiber сервер
func StartFiberServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.App.Port)
			log.Printf("🚀 Starting Fiber server on %s", addr)

			// Запускаем сервер в отдельной горутине
			go func() {
				if err := app.Listen(addr); err != nil {
					log.Fatalf("❌ Failed to start server: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 Shutting down Fiber server...")
			return app.Shutdown()
		},
	})
}

// NewQueries создает экземпляр generated.Queries из пула подключений
func NewQueries(pool *pgxpool.Pool) *generated.Queries {
	return generated.New(pool)
}

var AppModule = fx.Options(
	fx.Provide(
		NewFiberApp,
		NewPgPool,
		NewQueries,
	),
	RoleModule, // сюда входят все домены
	fx.Invoke(route.SetupRoutes),
	fx.Invoke(StartFiberServer),
)
