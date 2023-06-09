package routes

import (
	"github.com/dw-account-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initBalanceRoutes(router fiber.Router) {
	r := router.Group("/account/balance")

	r.Get("/inquiry/:uid", func(c *fiber.Ctx) error {
		return handlers.InquiryBalance(c)
	})

	// balance transaction
	// ---------------------------------------------------------------
	r.Post("/topup", func(c *fiber.Ctx) error {
		return handlers.TopupBalance(c)
	})

	r.Post("/deduct", func(c *fiber.Ctx) error {
		return handlers.DeductBalance(c)
	})

	// temporary commented out
	// ----------------------------------------
	//r.Post("/update", func(c *fiber.Ctx) error {
	//	return balances.UpdateBalance(c)
	//})

}
