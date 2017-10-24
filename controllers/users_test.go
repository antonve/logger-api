package controllers_test

import (
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

func init() {
	utils.SetupTesting()
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
		0,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	testUserToken, err := token.SignedString([]byte(config.GetConfig().JWTKey))

	assert.Nil(t, err)

	// Setup update request
	e := echo.New()
	updatedUser := strings.NewReader(`{
		"email": "new@example.com",
		"display_name": "OohNew",
		"password": "wotwotwot",
		"preferences": {
			"languages": ["JA", "KR", "ZH"],
			"public_profile": true
		}
  }`)
	req := httptest.NewRequest(echo.PUT, fmt.Sprintf("/api/user/%d", user.ID), updatedUser)

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
