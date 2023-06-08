package routes

import (
	"errors"
	"fmt"
	"github.com/dw-account-service/configs"
	"github.com/gofiber/fiber/v2"
)

func Start() error {
	config := configs.MainConfig.APIServer

	app := fiber.New()
	SetupRoutes(app)

	err := app.Listen(fmt.Sprintf(":%s", config.Port))
	if err != nil {
		return errors.New(fmt.Sprintf("error on starting service: %s", err.Error()))
	}

	return nil

}
