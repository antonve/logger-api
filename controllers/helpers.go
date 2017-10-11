package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/antonve/logger-api/config"
	"github.com/antonve/logger-api/models"

	"runtime/debug"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Serve a successful request
func Serve(context echo.Context, statusCode int) error {
	return context.JSONBlob(statusCode, []byte(`{"success": true}`))
}

// Serve a request with errors
func ServeWithError(context echo.Context, statusCode int, err error) error {
	handleError(err)
	body := []byte(fmt.Sprintf(`
		{
			"success": false,
			"errorCode": %d,
			"errorMessage": "%s"
		}`,
		statusCode,
		http.StatusText(statusCode)))
	return context.JSONBlob(statusCode, body)
}

// getUser helper
func getUser(context echo.Context) *models.User {
	token := context.Get("user")
	if token == nil {
		return nil
	}

	claims := token.(*jwt.Token).Claims
	if claims == nil {
		return nil
	}

	return claims.(*models.JwtClaims).User
}

func getRefreshTokenClaims(context echo.Context) *models.JwtRefreshTokenClaims {
	token := context.Get("user")
	if token == nil {
		return nil
	}

	claims := token.(*jwt.Token).Claims
	if claims == nil {
		return nil
	}

	return claims.(*models.JwtRefreshTokenClaims)
}

func handleError(err error) {
	log.Println(err.Error())

	if config.GetConfig().Debug {
		debug.PrintStack()
	}
}
