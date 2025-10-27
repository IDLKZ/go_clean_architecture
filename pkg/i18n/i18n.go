package i18n

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Поддерживаемые языки
const (
	LangRu = "ru" // Русский
	LangEn = "en" // Английский
	LangKk = "kk" // Казахский
)

var (
	// Bundle содержит все переводы
	Bundle *i18n.Bundle

	// DefaultLanguage - язык по умолчанию
	DefaultLanguage = LangRu

	// SupportedLanguages - список поддерживаемых языков
	SupportedLanguages = []string{LangRu, LangEn, LangKk}
)

//go:embed locales/*.json
var localesFS embed.FS

// Init инициализирует i18n bundle с переводами из встроенных файлов
func Init() error {
	// Создаем bundle с языком по умолчанию
	Bundle = i18n.NewBundle(language.Russian)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Загружаем переводы для каждого языка
	for _, lang := range SupportedLanguages {
		filename := fmt.Sprintf("locales/%s.json", lang)

		// Читаем файл из embed.FS
		data, err := localesFS.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read locale file %s: %w", filename, err)
		}

		// Загружаем сообщения напрямую из JSON
		if _, err := Bundle.ParseMessageFileBytes(data, filename); err != nil {
			return fmt.Errorf("failed to parse locale file %s: %w", filename, err)
		}
	}

	return nil
}

// GetLocalizer создает локализатор для указанного языка
// Если язык не поддерживается, используется язык по умолчанию
func GetLocalizer(lang string) *i18n.Localizer {
	if !IsLanguageSupported(lang) {
		lang = DefaultLanguage
	}
	return i18n.NewLocalizer(Bundle, lang)
}

// GetLocalizerFromAcceptLanguage создает локализатор на основе Accept-Language заголовка
func GetLocalizerFromAcceptLanguage(acceptLanguage string) *i18n.Localizer {
	return i18n.NewLocalizer(Bundle, acceptLanguage, DefaultLanguage)
}

// IsLanguageSupported проверяет, поддерживается ли указанный язык
func IsLanguageSupported(lang string) bool {
	for _, supported := range SupportedLanguages {
		if supported == lang {
			return true
		}
	}
	return false
}

// T переводит сообщение с указанным ID для заданного языка
// Поддерживает шаблонные переменные
func T(lang, messageID string, templateData map[string]interface{}) string {
	localizer := GetLocalizer(lang)

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

// TDefault переводит сообщение с дефолтным языком
func TDefault(messageID string, templateData map[string]interface{}) string {
	return T(DefaultLanguage, messageID, templateData)
}
