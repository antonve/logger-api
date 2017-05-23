package utils

import (
	"github.com/antonve/logger-api/controllers"

	"github.com/labstack/echo"
)

// SetupRouting Define all routes here
func SetupRouting(e *echo.Echo) {
	routesAPI := e.Group("/api")
	routesAPI.POST("/login", echo.HandlerFunc(controllers.APIUserLogin))
	routesAPI.POST("/register", echo.HandlerFunc(controllers.APIUserRegister))
}
