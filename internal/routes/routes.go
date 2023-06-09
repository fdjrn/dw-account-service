package routes

import (
	"errors"
	"fmt"
	"github.com/dw-account-service/configs"
	"github.com/gofiber/fiber/v2"
	"log"
)

func setupRoutes(app *fiber.App) {

	api := app.Group("/api/v1")
	initAccountRoutes(api)
	initBalanceRoutes(api)

	log.Println("[INIT] routes >> initialized")
}

func Initialize() error {
	config := configs.MainConfig.APIServer

	app := fiber.New()
	setupRoutes(app)

	err := app.Listen(fmt.Sprintf(":%s", config.Port))
	if err != nil {
		return errors.New(fmt.Sprintf("error on starting service: %s", err.Error()))
	}

	return nil

}
