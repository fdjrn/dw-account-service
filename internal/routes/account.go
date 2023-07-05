package routes

import (
	"github.com/dw-account-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initAccountRoutes(router fiber.Router) {
	accountHandler := handlers.NewAccountHandler()

	accountRoutes := router.Group("/account")

	accountRoutes.Post("/register", func(c *fiber.Ctx) error {
		return accountHandler.Register(c)
	})

	accountRoutes.Post("/unregister", func(c *fiber.Ctx) error {
		return accountHandler.Unregister(c)
	})

	accountRoutes.Post("/all", func(c *fiber.Ctx) error {
		return accountHandler.GetAccountsPaginated(c)
	})

	accountRoutes.Get("/:id", func(c *fiber.Ctx) error {
		return accountHandler.GetAccountByID(c)
	})

	accountRoutes.Post("/detail", func(c *fiber.Ctx) error {
		return accountHandler.GetAccount(c)
	})

	// -------------------------- Merchants --------------------------

	merchantRoutes := router.Group("/merchant")

	merchantRoutes.Post("/members", func(c *fiber.Ctx) error {
		return accountHandler.GetMerchantMembers(c, false)
	})

	merchantRoutes.Post("/members/period", func(c *fiber.Ctx) error {
		return accountHandler.GetMerchantMembers(c, true)
	})

	// -------------------------- TOOLS --------------------------

	accountRoutes.Post("/update-merchant-and-terminal", func(c *fiber.Ctx) error {
		return accountHandler.UpdateMerchantAndTerminalForAccount(c)
	})

	accountRoutes.Post("/sync-balance", func(c *fiber.Ctx) error {
		return accountHandler.SyncBalance(c)
	})

}
