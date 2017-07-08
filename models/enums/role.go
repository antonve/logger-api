package enums

import (
	"database/sql/driver"
	"errors"
)

// Role represents a user role
type (
	Role string
)

// Role values
const (
	RoleAdmin    Role = "ADMIN"
	RoleUser     Role = "USER"
	RoleDisabled Role = "DISABLED"
)

// Scan Role value
func (role *Role) Scan(src interface{}) error {
	if src == nil {
		return errors.New("This field cannot be NULL")
	}

	if stringRole, ok := src.([]byte); ok {
		*role = Role(string(stringRole[:]))

		return nil
	}

	return errors.New("Cannot convert enum to string")
}

// Value of Role
func (role Role) Value() (driver.Value, error) {
	return []byte(role), nil
}

// IsValid Role Value
func (role Role) IsValid() bool {
	if role == RoleAdmin {
		return true
	}
	if role == RoleUser {
		return true
	}
	if role == RoleDisabled {
		return true
	}

	return false
}
