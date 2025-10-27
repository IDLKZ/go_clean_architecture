package examples

import (
	i18nPkg "clean_architecture_fiber/pkg/i18n"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ExampleI18nBasicUsage демонстрирует базовое использование i18n
func ExampleI18nBasicUsage() {
	// Инициализация i18n (обычно происходит в main.go)
	if err := i18nPkg.Init(); err != nil {
		panic(err)
	}

	// Пример 1: Простой перевод на русский
	msgRu := i18nPkg.T("ru", "welcome", nil)
	fmt.Println("Russian:", msgRu) // Output: Добро пожаловать

	// Пример 2: Простой перевод на английский
	msgEn := i18nPkg.T("en", "welcome", nil)
	fmt.Println("English:", msgEn) // Output: Welcome

	// Пример 3: Простой перевод на казахский
	msgKk := i18nPkg.T("kk", "welcome", nil)
	fmt.Println("Kazakh:", msgKk) // Output: Қош келдіңіз
}

// ExampleI18nWithTemplateData демонстрирует использование шаблонных переменных
func ExampleI18nWithTemplateData() {
	if err := i18nPkg.Init(); err != nil {
		panic(err)
	}

	// Пример с шаблонными переменными
	templateData := map[string]interface{}{
		"Version": "2.0.0",
	}

	msgRu := i18nPkg.T("ru", "api.version", templateData)
	fmt.Println("Russian:", msgRu) // Output: API версия 2.0.0

	msgEn := i18nPkg.T("en", "api.version", templateData)
	fmt.Println("English:", msgEn) // Output: API version 2.0.0

	msgKk := i18nPkg.T("kk", "api.version", templateData)
	fmt.Println("Kazakh:", msgKk) // Output: API нұсқасы 2.0.0
}

// ExampleI18nInFiberHandler демонстрирует получение текущего языка в Fiber handler
func ExampleI18nInFiberHandler() {
	if err := i18nPkg.Init(); err != nil {
		panic(err)
	}

	app := fiber.New()

	// Подключаем i18n middleware
	app.Use(i18nPkg.Middleware())

	// Handler с получением текущего языка
	app.Get("/profile", func(c *fiber.Ctx) error {
		// Способ 1: Получить код текущего языка
		currentLang := i18nPkg.GetLanguage(c)
		fmt.Println("Current language:", currentLang) // Output: ru, en, или kk

		// Способ 2: Использовать перевод
		welcomeMsg := i18nPkg.Translate(c, "welcome", nil)

		// Способ 3: Перевод с шаблонными данными
		versionMsg := i18nPkg.Translate(c, "api.version", map[string]interface{}{
			"Version": "1.0.0",
		})

		return c.JSON(fiber.Map{
			"language": currentLang,
			"welcome":  welcomeMsg,
			"version":  versionMsg,
		})
	})

	// Пример использования в логике
	app.Get("/data", func(c *fiber.Ctx) error {
		currentLang := i18nPkg.GetLanguage(c)

		// Можно использовать язык для бизнес-логики
		var data interface{}
		switch currentLang {
		case i18nPkg.LangRu:
			data = "Русские данные"
		case i18nPkg.LangEn:
			data = "English data"
		case i18nPkg.LangKk:
			data = "Қазақша деректер"
		}

		return c.JSON(fiber.Map{
			"language": currentLang,
			"data":     data,
		})
	})
}

// ExampleI18nValidation демонстрирует использование i18n для валидации
func ExampleI18nValidation() {
	if err := i18nPkg.Init(); err != nil {
		panic(err)
	}

	// Пример с валидацией - обязательное поле
	requiredData := map[string]interface{}{
		"Field": "email",
	}

	msgRu := i18nPkg.T("ru", "validation.required", requiredData)
	fmt.Println("Russian:", msgRu) // Output: Поле email обязательно для заполнения

	msgEn := i18nPkg.T("en", "validation.required", requiredData)
	fmt.Println("English:", msgEn) // Output: Field email is required

	// Пример с валидацией - минимальная длина
	minLengthData := map[string]interface{}{
		"Field": "password",
		"Min":   8,
	}

	msgRu = i18nPkg.T("ru", "validation.min_length", minLengthData)
	fmt.Println("Russian:", msgRu) // Output: Поле password должно содержать минимум 8 символов

	msgEn = i18nPkg.T("en", "validation.min_length", minLengthData)
	fmt.Println("English:", msgEn) // Output: Field password must contain at least 8 characters
}

// ExampleI18nErrorMessages демонстрирует использование i18n для сообщений об ошибках
func ExampleI18nErrorMessages() {
	if err := i18nPkg.Init(); err != nil {
		panic(err)
	}

	// Различные ошибки на разных языках
	errors := []string{
		"error.not_found",
		"error.unauthorized",
		"error.forbidden",
		"error.internal_server",
	}

	for _, errorKey := range errors {
		fmt.Printf("\n%s:\n", errorKey)
		fmt.Printf("  RU: %s\n", i18nPkg.T("ru", errorKey, nil))
		fmt.Printf("  EN: %s\n", i18nPkg.T("en", errorKey, nil))
		fmt.Printf("  KK: %s\n", i18nPkg.T("kk", errorKey, nil))
	}
}

// ExampleI18nResourceMessages демонстрирует использование i18n для ресурсов (роли, разрешения)
func ExampleI18nResourceMessages() {
	if err := i18nPkg.Init(); err != nil {
		panic(err)
	}

	// Сообщения о ролях
	roleMessages := []string{
		"role_use_case.created",
		"role_use_case.updated",
		"role_use_case.deleted",
		"role_use_case.not_found",
	}

	fmt.Println("=== Role Messages ===")
	for _, msgKey := range roleMessages {
		fmt.Printf("\n%s:\n", msgKey)
		fmt.Printf("  RU: %s\n", i18nPkg.T("ru", msgKey, nil))
		fmt.Printf("  EN: %s\n", i18nPkg.T("en", msgKey, nil))
		fmt.Printf("  KK: %s\n", i18nPkg.T("kk", msgKey, nil))
	}

	// Сообщения о разрешениях
	permissionMessages := []string{
		"permission.created",
		"permission.updated",
		"permission.deleted",
		"permission.not_found",
	}

	fmt.Println("\n=== Permission Messages ===")
	for _, msgKey := range permissionMessages {
		fmt.Printf("\n%s:\n", msgKey)
		fmt.Printf("  RU: %s\n", i18nPkg.T("ru", msgKey, nil))
		fmt.Printf("  EN: %s\n", i18nPkg.T("en", msgKey, nil))
		fmt.Printf("  KK: %s\n", i18nPkg.T("kk", msgKey, nil))
	}
}
