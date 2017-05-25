package enums

import (
	"database/sql/driver"
	"errors"
)

// Language represents an language
type (
	Language string
)

// Language values
const (
	LanguageChinese  Language = "CN"
	LanguageJapanese Language = "JA"
	LanguageKorean   Language = "KR"
	LanguageMandarin Language = "ZH"
	LanguageGerman   Language = "DE"
)

// Scan Language value
func (language *Language) Scan(src interface{}) error {
	if src == nil {
		return errors.New("This field cannot be NULL")
	}

	if stringLanguage, ok := src.([]byte); ok {
		*language = Language(string(stringLanguage[:]))

		return nil
	}

	return errors.New("Cannot convert enum to string")
}

// Value of Language
func (language Language) Value() (driver.Value, error) {
	return []byte(language), nil
}

// IsValid Language Value
func (language Language) IsValid() bool {
	if language == LanguageJapanese {
		return true
	}
	if language == LanguageKorean {
		return true
	}
	if language == LanguageMandarin {
		return true
	}
	if language == LanguageGerman {
		return true
	}

	return false
}
