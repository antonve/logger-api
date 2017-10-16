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
	UserID   uint64 `json:"user_id" db:"user_id"`
	DeviceID string `json:"device_id" db:"device_id"`
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

func (refreshToken *RefreshToken) GenerateRefreshToken() error {
	// Generate new token
	jwtRefreshToken, err := refreshToken.GenerateRefreshTokenString()
	if err != nil {
		return err
	}

	// Create refresh token
	refreshTokenCollection := RefreshTokenCollection{RefreshTokens: make([]RefreshToken, 0)}
	_, err = refreshTokenCollection.Add(refreshToken)
	if err != nil {
		return err
	}

	refreshToken.RefreshToken = jwtRefreshToken

	return nil
}

// GenerateRefreshToken generates a new refresh  token that's valid for one year
// for a given user and device and returns the signed JWT token
func (refreshToken *RefreshToken) GenerateRefreshTokenString() (string, error) {
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

// Get a refresh token by claims
// nil is returned when a token is invalidated
func (refreshTokenCollection *RefreshTokenCollection) GetByClaims(claims *JwtRefreshTokenClaims) (*RefreshToken, error) {
	db := GetDatabase()
	defer db.Close()

	// Init refresh token
	refreshToken := RefreshToken{}

	// Get refresh token
	stmt, err := db.PrepareNamed(`
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
			user_id = :user_id AND
		  device_id = :device_id AND
			invalidated_at IS NULL
	`)
	if err != nil {
		return nil, err
	}

	stmt.Get(&refreshToken, claims)
	if refreshToken.ID == 0 {
		return nil, fmt.Errorf("no refresh token found with user id %v and device id %s", claims.UserID, claims.DeviceID)
	}

	return &refreshToken, nil
}

// Add a refresh token to the database
func (refreshTokenCollection *RefreshTokenCollection) Add(refreshToken *RefreshToken) (uint64, error) {
	db := GetDatabase()
	defer db.Close()

	// We must do the invalidation and creation of new tokens in a transaction
	// to make sure we don't leave the DB in a bad state if we crash
	tx, err := db.Beginx()

	// Invalidate older refresh tokens with the combination user_id, device_id
	invalidationQuery := `
		UPDATE refresh_tokens
		SET invalidated_at = NOW()
		WHERE
			user_id = :user_id AND
			device_id = :device_id AND
			invalidated_at IS NULL
	`
	_, err = tx.NamedExec(invalidationQuery, refreshToken)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert new token
	insertQuery := `
		INSERT INTO refresh_tokens (user_id, device_id, refresh_token)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err = tx.QueryRowx(insertQuery, refreshToken.UserID, refreshToken.DeviceID, refreshToken.RefreshToken).Scan(&refreshToken.ID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return refreshToken.ID, nil
}
