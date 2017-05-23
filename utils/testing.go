package utils

import (
	"github.com/antonve/logger-api/config"
	"github.com/antonve/logger-api/migrations"
	"fmt"
)

// SetupTesting the testing environment
func SetupTesting() {
	config.SetEnviroment(config.Environments["test"])

	teardown()

	err := migrations.Create()
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}

	migrations.Migrate()
}

func teardown() {
	err := migrations.Destroy()
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}
}
