package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/antonve/logger-api/config"
	"github.com/antonve/logger-api/controllers"
	"github.com/antonve/logger-api/models"
	"github.com/antonve/logger-api/models/enums"
	"github.com/antonve/logger-api/utils"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
)

type LoginBody struct {
	Token        string      `json:"token"`
	User         models.User `json:"user"`
	RefreshToken string      `json:"refresh_token"`
}

var mockSessionToken string
var mockSessionUser *models.User

func init() {
	utils.SetupTesting()
	mockSessionToken, mockSessionUser = utils.SetupTestUser("session_test")
}

func TestCreateUser(t *testing.T) {
	// Setup registration request
	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/api/register", strings.NewReader(`{"email": "register_test@example.com", "display_name": "logger", "password": "password"}`))
	assert.Nil(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controllers.APISessionRegister(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, `{"success": true}`, rec.Body.String())
	}
}

func TestCreateInvalidUser(t *testing.T) {
	// Setup registration request
	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/api/register", strings.NewReader(`{"email": "register_test@invalid##", "display_name": "invalid", "password": "password"}`))
	assert.Nil(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controllers.APISessionRegister(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotEqual(t, `{"success": true}`, rec.Body.String())
	}
}

func TestLoginUser(t *testing.T) {
	// Setup user to test login with
	user := models.User{Email: "login_test@example.com", DisplayName: "logger_user", Password: "password", Role: enums.RoleAdmin}
	user.HashPassword()
	userCollection := models.UserCollection{}
	userCollection.Add(&user)

	// Setup login request
	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/api/login", strings.NewReader(`{"email": "login_test@example.com", "password": "password", "device_id": "6db435f352d7ea4a67807a3feb447bf7"}`))
	assert.Nil(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controllers.APISessionLogin(c)) {
		// Check login response
		var body LoginBody
		assert.Equal(t, http.StatusOK, rec.Code)
		err = json.Unmarshal(rec.Body.Bytes(), &body)

		// Check if the user has information
		assert.Nil(t, err)
		assert.NotEmpty(t, body.Token)
		assert.NotEmpty(t, body.RefreshToken)
		assert.NotNil(t, body.User)

		// Check if the user has the correct information
		assert.Equal(t, "login_test@example.com", body.User.Email)
		assert.Equal(t, "logger_user", body.User.DisplayName)
		assert.Equal(t, enums.RoleAdmin, body.User.Role)

		// Make sure password is not sent back to the client
		assert.Empty(t, body.User.Password)
	}
}

func TestRefreshJWTToken(t *testing.T) {
	// Setup refresh request
	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/api/session/refresh", nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockSessionToken))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))(controllers.APISessionRefreshJWTToken)(c)) {
		// Check login response
		var body LoginBody
		assert.Equal(t, http.StatusOK, rec.Code)
		err = json.Unmarshal(rec.Body.Bytes(), &body)

		// Check if the user has information
		assert.Nil(t, err)
		assert.NotEmpty(t, body.Token)

		// Might want to check if the new token is usable
	}
}

func TestAuthenticateWithRefreshToken(t *testing.T) {
	// Setup refresh token
	refreshToken := models.RefreshToken{UserID: mockSessionUser.ID, DeviceID: "6db435f352d7ea4a67807a3feb447666"}
	jwtRefreshToken, err := refreshToken.GenerateRefreshTokenString()
	assert.Nil(t, err)
	refreshTokenCollection := models.RefreshTokenCollection{RefreshTokens: make([]models.RefreshToken, 0)}
	_, err = refreshTokenCollection.Add(&refreshToken)
	assert.Nil(t, err)

	// Setup authentication request
	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/api/session/authenticate", strings.NewReader(`{"device_id": "6db435f352d7ea4a67807a3feb447666"}`))
	assert.Nil(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtRefreshToken))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtRefreshTokenClaims{}))(controllers.APISessionAuthenticateWithRefreshToken)(c)) {
		// Check login response
		var body LoginBody
		assert.Equal(t, http.StatusOK, rec.Code)
		err = json.Unmarshal(rec.Body.Bytes(), &body)

		// Check if the user has information
		assert.Nil(t, err)
		assert.NotEmpty(t, body.Token)

		// Might want to check if the new token is usable
	}
}
