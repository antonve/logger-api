package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/antonve/logger-api/config"
	"github.com/antonve/logger-api/models"
	"github.com/antonve/logger-api/models/enums"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// APISessionLogin checks if user exists in database and returns jwt token if valid
func APISessionLogin(context echo.Context) error {
	// Attempt to bind request to User struct
	user := &models.User{}
	err := context.Bind(user)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	// Get authentication data
	userCollection := models.UserCollection{Users: make([]models.User, 0)}
	dbUser, err := userCollection.GetAuthenticationData(user.Email)
	if err != nil {
		log.Println(err)
		return echo.ErrUnauthorized
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		log.Println(err)
		return echo.ErrUnauthorized
	}

	// Set custom claims
	encodedToken, err := generateJWTToken(dbUser, 0)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return context.JSON(http.StatusOK, map[string]interface{}{
		"token": encodedToken,
		"user":  dbUser,
	})
}

// generateJWTToken generates a new JWT token that's valid for one hour for a given user
func generateJWTToken(user *models.User, refreshTokenID uint64) (string, error) {
	// Empty out password in case it was passed along
	user.Password = ""

	// Set claims
	claims := models.JwtClaims{
		user,
		1,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString([]byte(config.GetConfig().JWTKey))
}

// APISessionRefreshJWTToken will provide a new JWT token for a user who currently
// possess a valid JWT token
func APISessionRefreshJWTToken(context echo.Context) error {
	// Get user to work with
	user := getUser(context)
	if user == nil {
		return ServeWithError(context, 500, fmt.Errorf("could not receive user"))
	}

	// Get authentication data
	userCollection := models.UserCollection{Users: make([]models.User, 0)}
	dbUser, err := userCollection.GetAuthenticationData(user.Email)
	if err != nil {
		log.Println(err)
		return echo.ErrUnauthorized
	}

	// Check if refresh token was used for the current JWT token,
	// and if so make sure that one is still valid before refreshing
	refreshTokenStringID := context.Param("refresh_token_id")
	if refreshTokenStringID != "" {
		refreshTokenID, err := strconv.ParseUint(refreshTokenStringID, 10, 64)
		if err != nil {
			return ServeWithError(context, 500, err)
		}

		// Only proceed when we have a valid id
		if refreshTokenID != 0 {
			refreshTokenCollection := models.RefreshTokenCollection{RefreshTokens: make([]models.RefreshToken, 0)}
			refreshToken, err := refreshTokenCollection.Get(refreshTokenID)

			if err != nil {
				return ServeWithError(context, 500, err)
			}

			// Deny request if we have invalidated the refresh token
			if refreshToken.InvalidatedAt.Valid {
				log.Println("attempted JWT token refresh with expired session")
				return echo.ErrUnauthorized
			}
		}
	}

	// Set custom claims
	encodedToken, err := generateJWTToken(dbUser, 0)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	// Send new token to the user
	return context.JSON(http.StatusOK, map[string]interface{}{
		"token": encodedToken,
		"user":  dbUser,
	})
}

// APISessionNewJWTToken will provide a new JWT token for a user whose token has expired
// but can provide a device id & refresh token to generate a new one
//func APISessionNewJWTToken(context echo.Context) error {
//}

// APISessionRegister registers new user
func APISessionRegister(context echo.Context) error {
	user := &models.User{}

	// Attempt to bind request to User struct
	err := context.Bind(user)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	user.HashPassword()

	// Set default role
	user.Role = enums.RoleUser

	// Validate request
	err = user.Validate()
	if err != nil {
		return ServeWithError(context, 400, err)
	}

	// Save to database
	userCollection := models.UserCollection{}
	_, err = userCollection.Add(user)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return Serve(context, 201)
}
