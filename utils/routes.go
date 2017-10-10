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
	// Middleware
	authenticated := middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtClaims{}))
	//authenticatedWithRefreshToken := middleware.JWTWithConfig(config.GetJWTConfig(&models.JwtRefreshTokenClaims{}))

	// Routes
	routesAPI := e.Group("/api")
	routesAPI.POST("/login", echo.HandlerFunc(controllers.APISessionLogin))
	routesAPI.POST("/register", echo.HandlerFunc(controllers.APISessionRegister))

	routesSessions := routesAPI.Group("/session")
	routesSessions.POST("/refresh", authenticated(echo.HandlerFunc(controllers.APISessionRefreshJWTToken)))
	routesSessions.POST("/new", authenticated(echo.HandlerFunc(controllers.APISessionCreateRefreshToken)))
	//routesSessions.POST("/authenticate", authenticatedWithRefreshToken(echo.HandlerFunc(controllers.APISessionAuthenticateWithRefreshToken)))

	routesLogs := routesAPI.Group("/logs")
	routesLogs.Use(authenticated)
	routesLogs.GET("", echo.HandlerFunc(controllers.APILogsGetAll))
	routesLogs.POST("", echo.HandlerFunc(controllers.APILogsPost))
	routesLogs.GET("/:id", echo.HandlerFunc(controllers.APILogsGetByID))
	routesLogs.PUT("/:id", echo.HandlerFunc(controllers.APILogsUpdate))
	routesLogs.DELETE("/:id", echo.HandlerFunc(controllers.APILogsDelete))

}
