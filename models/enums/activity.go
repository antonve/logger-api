package enums

import (
	"database/sql/driver"
	"errors"
)

// Activity represents an activity that can be logged
type (
	Activity string
)

// Activity values
const (
	ActivityFlashcards  Activity = "FLASHCARDS"
	ActivityTextbook    Activity = "TEXTBOOK"
	ActivityReading     Activity = "READING"
	ActivityListening   Activity = "LISTENING"
	ActivityTranslation Activity = "TRANSLATION"
	ActivityGrammar     Activity = "GRAMMAR"
	ActivityOther       Activity = "OTHER"
)

// Scan Activity value
func (role *Activity) Scan(src interface{}) error {
	if src == nil {
		return errors.New("This field cannot be NULL")
	}

	if stringActivity, ok := src.([]byte); ok {
		*role = Activity(string(stringActivity[:]))

		return nil
	}

	return errors.New("Cannot convert enum to string")
}

// Value of Activity
func (role Activity) Value() (driver.Value, error) {
	return []byte(role), nil
}

// IsValid Activity Value
func (role Activity) IsValid() bool {
	if role == ActivityFlashcards {
		return true
	}
	if role == ActivityTextbook {
		return true
	}
	if role == ActivityReading {
		return true
	}
	if role == ActivityListening {
		return true
	}
	if role == ActivityTranslation {
		return true
	}
	if role == ActivityGrammar {
		return true
	}
	if role == ActivityOther {
		return true
	}

	return false
}
