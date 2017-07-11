package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/antonve/logger-api/config"
	"github.com/antonve/logger-api/controllers"
	"github.com/antonve/logger-api/models"
	"github.com/antonve/logger-api/models/enums"
	"github.com/antonve/logger-api/utils"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
)

type LoginBody struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func init() {
	utils.SetupTesting()
}

func TestUserCreateUser(t *testing.T) {
	// Setup registration request
	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/api/register", strings.NewReader(`{"email": "register_test@example.com", "display_name": "logger", "password": "password", "preferences": { "languages": ["JA", "DE", "ZH", "KR"], "public_profile": false }}`))
	if !assert.NoError(t, err) {
		return
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controllers.APIUserRegister(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, `{"success": true}`, rec.Body.String())
	}
}

//
// func TestUserCreateInvalidUser(t *testing.T) {
// 	// Setup registration request
// 	e := echo.New()
// 	req, err := http.NewRequest(echo.POST, "/api/register", strings.NewReader(`{"email": "register_test@invalid##", "display_name": "invalid", "password": "password"}`))
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
//
// 	if assert.NoError(t, controllers.APIUserRegister(c)) {
// 		assert.Equal(t, http.StatusBadRequest, rec.Code)
// 		assert.NotEqual(t, `{"success": true}`, rec.Body.String())
// 	}
// }

func TestUserLoginUser(t *testing.T) {
	// Setup user to test login with
	user := models.User{Email: "login_test@example.com", DisplayName: "logger_user", Password: "password", Role: enums.RoleAdmin}
	user.HashPassword()
	userCollection := models.UserCollection{}
	userCollection.Add(&user)

	// Setup login request
	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/api/login", strings.NewReader(`{"email": "login_test@example.com", "password": "password"}`))
	if !assert.NoError(t, err) {
		return
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controllers.APIUserLogin(c)) {
		// Check login response
		var body LoginBody
		assert.Equal(t, http.StatusOK, rec.Code)
		err = json.Unmarshal(rec.Body.Bytes(), &body)

		// Check if the user has information
		assert.Nil(t, err)
		assert.NotEmpty(t, body.Token)
		assert.NotNil(t, body.User)

		// Check if the user has the correct information
		assert.Equal(t, "login_test@example.com", body.User.Email)
		assert.Equal(t, "logger_user", body.User.DisplayName)
		assert.Equal(t, enums.RoleAdmin, body.User.Role)

		// Make sure password is not sent back to the client
		assert.Empty(t, body.User.Password)
	}
}

func TestUserUpdate(t *testing.T) {
	// Setup log to grab
	userCollection := models.UserCollection{}
	preferences := &models.Preferences{
		Languages:     []enums.Language{enums.LanguageGerman, enums.LanguageKorean},
		PublicProfile: false,
	}
	user := &models.User{
		Email:       "update_user@example.com",
		DisplayName: "UpdateMe",
		Password:    "examplepassword",
		Role:        enums.RoleUser,
		Preferences: *preferences,
	}
	user.ID, _ = userCollection.Add(user)
	user.Password = ""

	// Create JWT token with claims
	claims := models.JwtClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	testUserToken, err := token.SignedString([]byte(config.GetConfig().JWTKey))

	assert.Nil(t, err)

	// Setup log request
	e := echo.New()
	logUser := strings.NewReader(`{
		"email": "new@example.com",
		"display_name": "OohNew",
		"password": "wotwotwot",
		"preferences": {
			"languages": ["JA", "KR", "ZH"],
			"public_profile": true
		}
  }`)
	req := httptest.NewRequest(echo.PUT, fmt.Sprintf("/api/user/%d", user.ID), logUser)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testUserToken))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/user/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", user.ID))

	if assert.NoError(t, middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))(controllers.APIUserUpdate)(c)) {
		// Check response
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"success": true}`, rec.Body.String())

		updatedUser, _ := userCollection.Get(user.ID)
		assert.Equal(t, "new@example.com", updatedUser.Email)
		assert.Equal(t, "OohNew", updatedUser.DisplayName)
		assert.Equal(t, []enums.Language{enums.LanguageJapanese, enums.LanguageKorean, enums.LanguageMandarin}, updatedUser.Preferences.Languages)
	}
}
