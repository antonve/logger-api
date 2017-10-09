package utils

import (
	"github.com/antonve/logger-api/config"
	"github.com/antonve/logger-api/controllers"
	"github.com/antonve/logger-api/models"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// SetupRouting Define all routes here
func SetupRouting(e *echo.Echo) {
	authenticated := middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))

	routesAPI := e.Group("/api")
	routesAPI.POST("/login", echo.HandlerFunc(controllers.APIUserLogin))
	routesAPI.POST("/register", echo.HandlerFunc(controllers.APIUserRegister))

	routesSessions := routesAPI.Group("/session")
	routesSessions.POST("/refresh", authenticated(echo.HandlerFunc(controllers.APIUserRefreshJWTToken)))

	routesLogs := routesAPI.Group("/logs")
	routesLogs.Use(authenticated)
	routesLogs.GET("", echo.HandlerFunc(controllers.APILogsGetAll))
	routesLogs.POST("", echo.HandlerFunc(controllers.APILogsPost))
	routesLogs.GET("/:id", echo.HandlerFunc(controllers.APILogsGetByID))
	routesLogs.PUT("/:id", echo.HandlerFunc(controllers.APILogsUpdate))
	routesLogs.DELETE("/:id", echo.HandlerFunc(controllers.APILogsDelete))

}
