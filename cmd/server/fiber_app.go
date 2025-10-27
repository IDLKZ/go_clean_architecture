package server

import (
	"clean_architecture_fiber/config"
	i18nPkg "clean_architecture_fiber/pkg/i18n"
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SettleFiberApp настраивает и запускает Fiber веб-сервер
// Инициализирует middleware, маршруты и запускает сервер
func SettleFiberApp(cfg *config.Config, ctx context.Context, pool *pgxpool.Pool) {
	// Создаем экземпляр Fiber приложения с конфигурацией
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

	// Настраиваем маршруты API
	setupRoutes(app, ctx, pool)

	// Запускаем сервер
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
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

// setupRoutes настраивает маршруты API
func setupRoutes(app *fiber.App, ctx context.Context, pool *pgxpool.Pool) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		// Получаем текущий язык запроса
		currentLang := i18nPkg.GetLanguage(c)

		// Пример использования i18n
		welcomeMsg := i18nPkg.Translate(c, "welcome", nil)

		return c.JSON(fiber.Map{
			"status":   "ok",
			"health":   "healthy",
			"language": currentLang,
			"message":  welcomeMsg,
		})
	})

	// API группа - все маршруты API начинаются с /api
	api := app.Group("/api")

	// Версия API v1
	v1 := api.Group("/v1")

	// TODO: Здесь будут добавлены маршруты для:
	// - v1.Group("/roles")         // Управление ролями
	// - v1.Group("/permissions")   // Управление разрешениями
	// - v1.Group("/users")         // Управление пользователями
	// - v1.Group("/auth")          // Аутентификация и авторизация

	// Временный тестовый маршрут с демонстрацией i18n
	v1.Get("/", func(c *fiber.Ctx) error {
		// Получаем текущий язык запроса
		currentLang := i18nPkg.GetLanguage(c)

		// Пример использования i18n с шаблонными переменными
		versionMsg := i18nPkg.Translate(c, "api.version", map[string]interface{}{
			"Version": "1.0.0",
		})

		return c.JSON(fiber.Map{
			"language": currentLang,
			"message":  versionMsg,
			"version":  "1.0.0",
		})
	})

	// 404 handler с поддержкой i18n
	app.Use(func(c *fiber.Ctx) error {
		errorMsg := i18nPkg.Translate(c, "error.not_found", nil)

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": errorMsg,
			"path":    c.Path(),
		})
	})
}
