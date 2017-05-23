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

	routesLogs := routesAPI.Group("/logs")
	// routesLogs.Use(middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{})))
	routesLogs.GET("", echo.HandlerFunc(controllers.APILogsGetAll))
	routesLogs.POST("", echo.HandlerFunc(controllers.APILogsCreate))
	routesLogs.GET("/:id", echo.HandlerFunc(controllers.APILogsGetByID))
	routesLogs.PUT("/:id", echo.HandlerFunc(controllers.APILogsUpdate))
	routesLogs.DELETE("/:id", echo.HandlerFunc(controllers.APILogsDelete))

}
