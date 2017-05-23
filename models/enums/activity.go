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
func (activity *Activity) Scan(src interface{}) error {
	if src == nil {
		return errors.New("This field cannot be NULL")
	}

	if stringActivity, ok := src.([]byte); ok {
		*activity = Activity(string(stringActivity[:]))

		return nil
	}

	return errors.New("Cannot convert enum to string")
}

// Value of Activity
func (activity Activity) Value() (driver.Value, error) {
	return []byte(activity), nil
}

// IsValid Activity Value
func (activity Activity) IsValid() bool {
	if activity == ActivityFlashcards {
		return true
	}
	if activity == ActivityTextbook {
		return true
	}
	if activity == ActivityReading {
		return true
	}
	if activity == ActivityListening {
		return true
	}
	if activity == ActivityTranslation {
		return true
	}
	if activity == ActivityGrammar {
		return true
	}
	if activity == ActivityOther {
		return true
	}

	return false
}
