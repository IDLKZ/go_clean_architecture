package mapper

import (
	i18nPkg "clean_architecture_fiber/pkg/i18n"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// getLocalizedText выбирает локализованный текст на основе текущего языка запроса
// Если перевод для запрошенного языка отсутствует (Valid = false), используется русский (fallback)
//
// Параметры:
//   - ctx: Fiber контекст для определения текущего языка
//   - textRu: Текст на русском языке (обязательный, используется как fallback)
//   - textEn: Текст на английском языке (опциональный)
//   - textKk: Текст на казахском языке (опциональный)
//
// Возвращает: локализованный текст на основе языка запроса
func getLocalizedText(ctx *fiber.Ctx, textRu string, textEn pgtype.Text, textKk pgtype.Text) string {
	currentLanguage := i18nPkg.GetLanguage(ctx)

	switch currentLanguage {
	case i18nPkg.LangEn:
		// Английский язык
		if textEn.Valid {
			return textEn.String
		}
		return textRu

	case i18nPkg.LangKk:
		// Казахский язык
		if textKk.Valid {
			return textKk.String
		}
		return textRu

	default:
		// Русский язык (по умолчанию)
		return textRu
	}
}

// uuidToString безопасно преобразует pgtype.UUID в строку
// Если UUID невалидный, возвращает пустую строку
func uuidToString(uuid pgtype.UUID) string {
	if !uuid.Valid {
		return ""
	}

	uuidBytes, err := uuid.Value()
	if err != nil || uuidBytes == nil {
		return ""
	}

	uuidStr, ok := uuidBytes.(string)
	if !ok {
		return ""
	}

	return uuidStr
}
