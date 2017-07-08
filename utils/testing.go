package utils

import (
	"fmt"
	"time"

	"github.com/antonve/logger-api/config"
	"github.com/antonve/logger-api/migrations"
	"github.com/antonve/logger-api/models"
	"github.com/antonve/logger-api/models/enums"
	jwt "github.com/dgrijalva/jwt-go"
)

// SetupTesting the testing environment
func SetupTesting() {
	config.SetEnviroment(config.Environments["test"])

	teardown()

	err := migrations.Create()
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}

	migrations.Migrate()
}

// SetupTestUser a mock user for testing
func SetupTestUser() (string, *models.User) {
	user := &models.User{
		Username:    "mock_user",
		DisplayName: "mock_name",
		Password:    "mock_password",
		Role:        enums.RoleUser,
	}

	// Get authentication data
	userCollection := models.UserCollection{Users: make([]models.User, 0)}
	_, err := userCollection.Add(user)
	if err != nil {
		return "", nil
	}

	dbUser, err := userCollection.GetAuthenticationData(user.Username)
	if err != nil {
		return "", nil
	}

	// Set custom claims
	dbUser.Password = ""
	claims := models.JwtClaims{
		dbUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	encodedToken, err := token.SignedString([]byte(config.GetConfig().JWTKey))
	if err != nil {
		return "", nil
	}

	return encodedToken, dbUser
}

func teardown() {
	err := migrations.Destroy()
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}
}
