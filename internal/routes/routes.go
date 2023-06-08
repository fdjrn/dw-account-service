package routes

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api/v1")
	initAccountRoutes(api)
	initBalanceRoutes(api)

	log.Println("[INIT] routes >> initialized")
}
