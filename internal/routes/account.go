package routes

import (
	"github.com/dw-account-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initAccountRoutes(router fiber.Router) {
	r := router.Group("/account")

	r.Post("/register", func(c *fiber.Ctx) error {
		return handlers.Register(c)
	})

	r.Post("/unregister", func(c *fiber.Ctx) error {
		return handlers.Unregister(c)
	})

	// ---------------------------------------------------------------------------------------------------------------
	// it can use query params to filter their active status
	// example:
	// api/v1/account?active=true 	--> to fetch only active account
	// api/v1/account?active=false --> to fetch only unregistered account
	// api/v1/account 				--> to fetch all registered account whether its active or unregistered
	// ---------------------------------------------------------------------------------------------------------------
	// -- DEPRECATED --
	// ---------------------------------------------------------------------------------------------------------------
	r.Get("", func(c *fiber.Ctx) error {
		//return handlers.GetAllRegisteredAccount(c)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "-- deprecated --",
			"data":    nil,
		})
	})

	r.Post("/all", func(c *fiber.Ctx) error {
		return handlers.GetAllRegisteredAccountPaginated(c)
	})

	r.Get("/:id", func(c *fiber.Ctx) error {
		return handlers.GetRegisteredAccount(c)
	})

	r.Get("/uid/:uid", func(c *fiber.Ctx) error {
		return handlers.GetRegisteredAccountByUID(c)
	})

}
