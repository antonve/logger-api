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

// APIUserLogin checks if user exists in database and returns jwt token if valid
func APIUserLogin(context echo.Context) error {
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
	encodedToken, err := generateJWTToken(dbUser)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return context.JSON(http.StatusOK, map[string]interface{}{
		"token": encodedToken,
		"user":  dbUser,
	})
}

// generateJWTToken generates a new JWT token that's valid for one hour for a given user
func generateJWTToken(user *models.User) (string, error) {
	user.Password = ""
	claims := models.JwtClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString([]byte(config.GetConfig().JWTKey))
}

// APIUserRefreshJWTToken will provide a new JWT token for a user who currently
// possess a valid JWT token
func APIUserRefreshJWTToken(context echo.Context) error {
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

	// Set custom claims
	encodedToken, err := generateJWTToken(dbUser)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	// Send new token to the user
	return context.JSON(http.StatusOK, map[string]interface{}{
		"token": encodedToken,
		"user":  dbUser,
	})
}

// APIUserNewJWTToken will provide a new JWT token for a user whose token has expired
// but can provide a device id & refresh token to generate a new one
//func APIUserNewJWTToken(context echo.Context) error {
//}

// APIUserRegister registers new user
func APIUserRegister(context echo.Context) error {
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

// APIUserGetAll gets all users
func APIUserGetAll(context echo.Context) error {
	userCollection := models.UserCollection{Users: make([]models.User, 0)}
	err := userCollection.GetAll()

	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return context.JSON(http.StatusOK, userCollection)
}

// APIUserGetByID get the profile of a user
func APIUserGetByID(context echo.Context) error {
	userCollection := models.UserCollection{}

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	user, err := userCollection.Get(id)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	if user == nil {
		return ServeWithError(context, 404, fmt.Errorf("No User found with id %v", id))
	}

	return context.JSON(http.StatusOK, user)
}

// APIUserUpdate updates a user
func APIUserUpdate(context echo.Context) error {
	user := &models.User{}

	// Attempt to bind request to User struct
	err := context.Bind(user)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	// Parse out id
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return ServeWithError(context, 500, err)
	}
	user.ID = id

	// Update
	userCollection := models.UserCollection{}
	err = userCollection.Update(user)
	if err != nil {
		return ServeWithError(context, 500, err)
	}

	return Serve(context, 200)
}
