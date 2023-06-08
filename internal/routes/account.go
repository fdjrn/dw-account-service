package routes

import (
	"github.com/dw-account-service/internal/handlers/account"
	"github.com/gofiber/fiber/v2"
)

func initAccountRoutes(router fiber.Router) {
	r := router.Group("/account")

	r.Post("/register", func(c *fiber.Ctx) error {
		return account.Register(c)
	})

	r.Post("/unregister", func(c *fiber.Ctx) error {
		return account.Unregister(c)
	})

	// it can use query params to filter their active status
	// example:
	//
	//	api/v1/account?active=true 	--> to fetch only active account
	//	api/v1/account?active=false --> to fetch only unregistered account
	//	api/v1/account 				--> to fetch all registered account whether its active or unregistered
	r.Get("", func(c *fiber.Ctx) error {
		return account.GetAllRegisteredAccount(c)
	})

	r.Get("/:id", func(c *fiber.Ctx) error {
		return account.GetRegisteredAccount(c)
	})

	r.Get("/uid/:uid", func(c *fiber.Ctx) error {
		return account.GetRegisteredAccountByUID(c)
	})

}
