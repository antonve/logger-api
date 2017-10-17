package models

import (
	"errors"
	"fmt"

	"github.com/antonve/logger-api/models/enums"
	"github.com/badoux/checkmail"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
)

// UserCollection array of users
type UserCollection struct {
	Users []User `json:"users"`
}

// User model
type User struct {
	ID          uint64     `json:"id" db:"id"`
	Email       string     `json:"email" db:"email"`
	DisplayName string     `json:"display_name" db:"display_name"`
	Password    string     `json:"password" db:"password"`
	Role        enums.Role `json:"role" db:"role"`
}

// JwtClaims json web token claim
type JwtClaims struct {
	User           *User  `json:"user"`
	RefreshTokenID uint64 `json:"refresh_token_id"`
	jwt.StandardClaims
}

// Length returns the amount of users in the collection
func (userCollection *UserCollection) Length() int {
	return len(userCollection.Users)
}

// Validate the User model
func (user *User) Validate() error {
	if len(user.Email) == 0 {
		return errors.New("no `email` supplied")
	}
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("invalid `Email` supplied")
	}
	if len(user.DisplayName) == 0 {
		return errors.New("invalid `DisplayName` supplied")
	}
	if len(user.Role) == 0 || !user.Role.IsValid() {
		return errors.New("invalid `Role` supplied")
	}
	if user.ID == 0 && len(user.Password) == 0 {
		return errors.New("invalid `Password` supplied")
	}

	return nil
}

// HashPassword hash the currently set password
func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return nil
}

// GetAll returns all users
func (userCollection *UserCollection) GetAll() error {
	db := GetDatabase()
	defer db.Close()

	err := db.Select(&userCollection.Users, `
		SELECT
			id,
			email,
			display_name,
			role
		FROM users
	`)

	return err
}

// Get a user by id
func (userCollection *UserCollection) Get(id uint64) (*User, error) {
	db := GetDatabase()
	defer db.Close()

	// Init user
	user := User{}

	// Get user
	stmt, err := db.Preparex(`
		SELECT
			id,
			email,
			display_name,
			role
		FROM users
		WHERE
			id = $1
	`)
	if err != nil {
		return nil, err
	}

	stmt.Get(&user, id)
	return &user, nil
}

// GetAuthenticationData get data needed to generate jwt token
func (userCollection *UserCollection) GetAuthenticationData(email string) (*User, error) {
	db := GetDatabase()
	defer db.Close()

	user := User{}

	stmt, err := db.Preparex(`
		SELECT
			id,
			email,
			display_name,
			role,
			password
		FROM users
		WHERE email = $1
	`)
	if err != nil {
		return nil, err
	}

	stmt.Get(&user, email)

	return &user, err
}

// Add a user to the database
func (userCollection *UserCollection) Add(user *User) (uint64, error) {
	db := GetDatabase()
	defer db.Close()

	query := `
		INSERT INTO users
		(email, display_name, password, role)
		VALUES (:email, :display_name, :password, :role)
		RETURNING id
	`
	rows, err := db.NamedQuery(query, user)

	if err != nil {
		return 0, err
	}

	var id uint64
	if rows.Next() {
		rows.Scan(&id)
	}

	return id, nil
}

// Update a user
func (userCollection *UserCollection) Update(user *User) error {
	db := GetDatabase()
	defer db.Close()

	query := `
		UPDATE users
		SET
			email = :email,
			display_name = :display_name,
			role = :role
		WHERE id = :id
	`
	result, err := db.NamedExec(query, user)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if rows == 0 {
		err = fmt.Errorf("No user found with id %v", user.ID)
	}

	return err
}
