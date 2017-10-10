package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/antonve/logger-api/config"
	"github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwtRefreshTokenClaims json web token claim for refresh token
type JwtRefreshTokenClaims struct {
	UserID   uint64 `json:"user_id"`
	DeviceID string `json:"device_id"`
	jwt.StandardClaims
}

// RefreshTokenCollection array of refresh tokens
type RefreshTokenCollection struct {
	RefreshTokens []RefreshToken `json:"refresh_tokens"`
}

// RefreshToken model
type RefreshToken struct {
	ID            uint64      `json:"id" db:"id"`
	UserID        uint64      `json:"user_id" db:"user_id"`
	DeviceID      string      `json:"device_id" db:"device_id"`
	RefreshToken  string      `json:"refresh_token" db:"refresh_token"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	InvalidatedAt pq.NullTime `json:"invalidated_at" db:"invalidated_at"`
}

// GenerateRefreshToken generates a new refresh  token that's valid for one year
// for a given user and device and returns the signed JWT token
func (refreshToken *RefreshToken) GenerateRefreshToken() (string, error) {
	// Set claims
	claims := JwtRefreshTokenClaims{
		refreshToken.UserID,
		refreshToken.DeviceID,
		jwt.StandardClaims{
			// Duration of 1 year
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.GetConfig().JWTKey))
	if err != nil {
		return "", err
	}

	// Hash signed token to store in DB
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(signedToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	refreshToken.RefreshToken = string(hashedToken)

	return signedToken, nil
}

// Length returns the amount of refresh tokens in the collection
func (refreshTokenCollection *RefreshTokenCollection) Length() int {
	return len(refreshTokenCollection.RefreshTokens)
}

// Validate the RefreshToken model
func (refreshToken *RefreshToken) Validate() error {
	if refreshToken.UserID == 0 {
		return errors.New("invalid `UserID` supplied")
	}
	if refreshToken.DeviceID == "" {
		return errors.New("invalid `DeviceID` supplied")
	}
	if refreshToken.CreatedAt.IsZero() {
		return errors.New("invalid `CreatedAt` supplied")
	}
	if refreshToken.UpdatedAt.IsZero() {
		return errors.New("invalid `UpdatedAt` supplied")
	}

	return nil
}

// Get a refresh token by id
func (refreshTokenCollection *RefreshTokenCollection) Get(id uint64) (*RefreshToken, error) {
	db := GetDatabase()
	defer db.Close()

	// Init refresh token
	refreshToken := RefreshToken{}

	// Get refresh token
	stmt, err := db.Preparex(`
		SELECT
			id,
			user_id,
			device_id,
			refresh_token,
			created_at,
			updated_at,
			invalidated_at
		FROM refresh_tokens
		WHERE
			id = $1
	`)
	if err != nil {
		return nil, err
	}

	stmt.Get(&refreshToken, id)
	if refreshToken.ID == 0 {
		return nil, fmt.Errorf("no refresh token found with id %v", id)
	}

	return &refreshToken, nil
}

// Add a refresh token to the database
func (refreshTokenCollection *RefreshTokenCollection) Add(refreshToken *RefreshToken) (uint64, error) {
	db := GetDatabase()
	defer db.Close()

	query := `
		INSERT INTO refresh_tokens (user_id, device_id, refresh_token)
		VALUES (:user_id, :device_id, :refresh_token)
		RETURNING id
	`
	rows, err := db.NamedQuery(query, refreshToken)

	if err != nil {
		return 0, err
	}

	var id uint64
	if rows.Next() {
		rows.Scan(&id)
	}

	return id, nil
}
