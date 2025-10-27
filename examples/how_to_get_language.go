package examples

import (
	i18nPkg "clean_architecture_fiber/pkg/i18n"
	"github.com/gofiber/fiber/v2"
)

// HowToGetLanguage демонстрирует различные способы получения текущего языка в проекте
func HowToGetLanguage() {
	// Инициализация i18n (обычно происходит в main.go)
	if err := i18nPkg.Init(); err != nil {
		panic(err)
	}

	app := fiber.New()

	// Подключаем i18n middleware (обязательно!)
	app.Use(i18nPkg.Middleware())

	// ============================================================================
	// СПОСОБ 1: Использование i18nPkg.GetLanguage(c)
	// ============================================================================
	// Самый простой и рекомендуемый способ
	app.Get("/method1", func(c *fiber.Ctx) error {
		// Получаем текущий язык из контекста
		lang := i18nPkg.GetLanguage(c)
		// lang будет: "ru", "en", или "kk"

		return c.JSON(fiber.Map{
			"method":   "GetLanguage",
			"language": lang,
		})
	})

	// ============================================================================
	// СПОСОБ 2: Использование через c.Locals напрямую
	// ============================================================================
	// Более низкоуровневый способ, если нужен полный контроль
	app.Get("/method2", func(c *fiber.Ctx) error {
		// Получаем значение напрямую из Fiber контекста
		lang, ok := c.Locals(i18nPkg.LanguageContextKey).(string)
		if !ok || lang == "" {
			lang = i18nPkg.DefaultLanguage
		}

		return c.JSON(fiber.Map{
			"method":   "Locals direct access",
			"language": lang,
		})
	})

	// ============================================================================
	// СПОСОБ 3: Использование константы для проверки конкретного языка
	// ============================================================================
	app.Get("/method3", func(c *fiber.Ctx) error {
		lang := i18nPkg.GetLanguage(c)

		// Можно сравнивать с константами
		var message string
		switch lang {
		case i18nPkg.LangRu:
			message = "Вы используете русский язык"
		case i18nPkg.LangEn:
			message = "You are using English"
		case i18nPkg.LangKk:
			message = "Сіз қазақ тілін пайдаланасыз"
		default:
			message = "Unknown language"
		}

		return c.JSON(fiber.Map{
			"method":   "Language constants",
			"language": lang,
			"message":  message,
		})
	})

	// ============================================================================
	// СПОСОБ 4: Получение языка вместе с переводом
	// ============================================================================
	app.Get("/method4", func(c *fiber.Ctx) error {
		// Получаем язык
		lang := i18nPkg.GetLanguage(c)

		// И сразу используем для перевода
		welcomeMsg := i18nPkg.Translate(c, "welcome", nil)

		return c.JSON(fiber.Map{
			"method":   "Language with translation",
			"language": lang,
			"message":  welcomeMsg,
		})
	})

	// ============================================================================
	// СПОСОБ 5: Использование в middleware/business logic
	// ============================================================================
	// Можно использовать в собственных middleware
	app.Use(func(c *fiber.Ctx) error {
		lang := i18nPkg.GetLanguage(c)

		// Логика на основе языка
		if lang == i18nPkg.LangRu {
			// Специфическая логика для русского языка
			c.Set("X-Content-Language", "ru-RU")
		}

		return c.Next()
	})

	// ============================================================================
	// СПОСОБ 6: Проверка поддерживаемого языка
	// ============================================================================
	app.Get("/method6", func(c *fiber.Ctx) error {
		lang := i18nPkg.GetLanguage(c)

		// Проверяем, поддерживается ли язык
		isSupported := i18nPkg.IsLanguageSupported(lang)

		// Получаем список всех поддерживаемых языков
		supportedLangs := i18nPkg.SupportedLanguages

		return c.JSON(fiber.Map{
			"method":              "Language validation",
			"current_language":    lang,
			"is_supported":        isSupported,
			"supported_languages": supportedLangs,
		})
	})

	// ============================================================================
	// ПРИМЕЧАНИЯ:
	// ============================================================================
	// 1. Язык определяется в следующем порядке приоритета:
	//    - Query параметр: ?lang=en
	//    - HTTP заголовок: Accept-Language: en
	//    - Язык по умолчанию: ru
	//
	// 2. Middleware i18nPkg.Middleware() должен быть подключен ДО маршрутов
	//
	// 3. Поддерживаемые языки:
	//    - i18nPkg.LangRu = "ru" (Русский)
	//    - i18nPkg.LangEn = "en" (Английский)
	//    - i18nPkg.LangKk = "kk" (Казахский)
	//
	// 4. Если указан неподдерживаемый язык, используется DefaultLanguage (ru)
}
