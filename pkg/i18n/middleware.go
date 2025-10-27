package i18n

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// LocalizerContextKey - ключ для сохранения локализатора в контексте Fiber
const LocalizerContextKey = "localizer"

// LanguageContextKey - ключ для сохранения кода языка в контексте Fiber
const LanguageContextKey = "language"

// Middleware создает middleware для определения языка запроса
// Язык определяется в следующем порядке:
// 1. Query параметр ?lang=ru
// 2. HTTP заголовок Accept-Language
// 3. Язык по умолчанию (русский)
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Пытаемся получить язык из query параметра
		lang := c.Query("lang", "")

		var localizer *i18n.Localizer
		var detectedLang string

		if lang != "" && IsLanguageSupported(lang) {
			// Используем язык из query параметра
			localizer = GetLocalizer(lang)
			detectedLang = lang
		} else {
			// Используем Accept-Language заголовок
			acceptLanguage := c.Get("Accept-Language", "")
			if acceptLanguage == "" {
				detectedLang = DefaultLanguage
			} else {
				// Парсим Accept-Language и определяем язык
				detectedLang = parseAcceptLanguage(acceptLanguage)
			}
			localizer = GetLocalizerFromAcceptLanguage(acceptLanguage)
		}

		// Сохраняем локализатор и код языка в контексте запроса
		c.Locals(LocalizerContextKey, localizer)
		c.Locals(LanguageContextKey, detectedLang)

		return c.Next()
	}
}

// GetLocalizerFromContext извлекает локализатор из контекста Fiber
// Если локализатор не найден, возвращает локализатор с языком по умолчанию
func GetLocalizerFromContext(c *fiber.Ctx) *i18n.Localizer {
	localizer, ok := c.Locals(LocalizerContextKey).(*i18n.Localizer)
	if !ok {
		return GetLocalizer(DefaultLanguage)
	}
	return localizer
}

// Translate переводит сообщение с использованием локализатора из контекста
func Translate(c *fiber.Ctx, messageID string, templateData map[string]interface{}) string {
	localizer := GetLocalizerFromContext(c)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		// Если перевод не найден, возвращаем messageID
		return messageID
	}

	return msg
}

// MustTranslate переводит сообщение или возвращает defaultMessage при ошибке
func MustTranslate(c *fiber.Ctx, messageID string, defaultMessage string, templateData map[string]interface{}) string {
	localizer := GetLocalizerFromContext(c)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		return defaultMessage
	}

	return msg
}

// GetLanguage возвращает текущий код языка из контекста Fiber
// Если язык не найден в контексте, возвращает язык по умолчанию
//
// Возвращаемые значения: "ru", "en", "kk"
//
// Пример использования:
//
//	func handler(c *fiber.Ctx) error {
//	    lang := i18nPkg.GetLanguage(c)
//	    // lang будет одним из: "ru", "en", "kk"
//
//	    // Использование в логике
//	    if lang == i18nPkg.LangRu {
//	        // Логика для русского языка
//	    }
//
//	    return c.JSON(fiber.Map{"language": lang})
//	}
func GetLanguage(c *fiber.Ctx) string {
	lang, ok := c.Locals(LanguageContextKey).(string)
	if !ok || lang == "" {
		return DefaultLanguage
	}
	return lang
}

// parseAcceptLanguage парсит заголовок Accept-Language и возвращает первый поддерживаемый язык
// Если поддерживаемый язык не найден, возвращает язык по умолчанию
func parseAcceptLanguage(acceptLanguage string) string {
	if acceptLanguage == "" {
		return DefaultLanguage
	}

	// Простой парсинг: берем первые два символа (код языка)
	// Например: "en-US,en;q=0.9,ru;q=0.8" -> "en"
	if len(acceptLanguage) >= 2 {
		lang := acceptLanguage[:2]
		if IsLanguageSupported(lang) {
			return lang
		}
	}

	return DefaultLanguage
}
