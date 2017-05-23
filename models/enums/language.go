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
	LanguageJapanese Language = "JA"
	LanguageKorean   Language = "KR"
	LanguageMandarin Language = "ZH"
	LanguageGerman   Language = "DE"
)

// Scan Language value
func (role *Language) Scan(src interface{}) error {
	if src == nil {
		return errors.New("This field cannot be NULL")
	}

	if stringLanguage, ok := src.([]byte); ok {
		*role = Language(string(stringLanguage[:]))

		return nil
	}

	return errors.New("Cannot convert enum to string")
}

// Value of Language
func (role Language) Value() (driver.Value, error) {
	return []byte(role), nil
}

// IsValid Language Value
func (role Language) IsValid() bool {
	if role == LanguageJapanese {
		return true
	}
	if role == LanguageKorean {
		return true
	}
	if role == LanguageMandarin {
		return true
	}
	if role == LanguageGerman {
		return true
	}

	return false
}
