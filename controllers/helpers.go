package controllers

import (
	"log"
	"net/http"

	"github.com/antonve/logger-api/config"
	"github.com/antonve/logger-api/models"

	"runtime/debug"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Return201 helper
func Return201(context echo.Context) error {
	return context.JSONBlob(http.StatusCreated, []byte(`{"success": true}`))
}

// Return200 helper
func Return200(context echo.Context) error {
	return context.JSONBlob(http.StatusOK, []byte(`{"success": true}`))
}

// Return400 helper
func Return400(context echo.Context, err error) error {
	handleError(err)
	return Serve400(context)
}

// Serve400 helper
func Serve400(context echo.Context) error {
	return context.JSONBlob(http.StatusBadRequest, []byte(`{"success": false, "errorCode": 400, "errorMessage": "400 bad request"}`))
}

// Return403 helper
func Return403(context echo.Context, err error) error {
	handleError(err)
	return Serve403(context)
}

// Serve403 helper
func Serve403(context echo.Context) error {
	return context.JSONBlob(http.StatusForbidden, []byte(`{"success": false, "errorCode": 403, "errorMessage": "400 forbidden"}`))
}

// Return404 helper
func Return404(context echo.Context, err error) error {
	handleError(err)
	return Serve404(context)
}

// Serve404 helper
func Serve404(context echo.Context) error {
	return context.JSONBlob(http.StatusNotFound, []byte(`{"success": false, "errorCode": 404, "errorMessage": "404 page not found"}`))
}

// Serve405 helper
func Serve405(context echo.Context) error {
	return context.JSONBlob(http.StatusMethodNotAllowed, []byte(`{"success": false, "errorCode": 405, "errorMessage": "405 method not allowed"}`))
}

// Return500 helper
func Return500(context echo.Context, err error) error {
	handleError(err)
	return Serve500(context)
}

// Serve500 helper
func Serve500(context echo.Context) error {
	return context.JSONBlob(http.StatusInternalServerError, []byte(`{"success": false, "errorCode": 500, "errorMessage": "500 internal server error"}`))
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

func handleError(err error) {
	log.Println(err.Error())

	if config.GetConfig().Debug {
		debug.PrintStack()
	}
}
