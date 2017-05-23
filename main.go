package main

import (
	"github.com/antonve/logger-api/utils"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"log"
)

func main() {
	// Echo instance
	e := echo.New()
	log.Println("Starting Logger API")

	// Middleware
	e.Use(middleware.Recover())
	defer utils.SetupErrorLogging(e)()

	// Serve static assets
	utils.SetupStaticAssets(e)

	// Routes
	utils.SetupRouting(e)

	// Start server
	e.Start(":3000")
}
