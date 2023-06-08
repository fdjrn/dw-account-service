package routes

import (
	"github.com/dw-account-service/internal/handlers/balances"
	"github.com/gofiber/fiber/v2"
)

func initBalanceRoutes(router fiber.Router) {
	r := router.Group("/account/balance")

	r.Get("/inquiry/:uid", func(c *fiber.Ctx) error {
		return balances.InquiryBalance(c)
	})

	// balance transaction
	// ---------------------------------------------------------------
	r.Post("/topup", func(c *fiber.Ctx) error {
		return balances.TopupBalance(c)
	})

	r.Post("/deduct", func(c *fiber.Ctx) error {
		return balances.DeductBalance(c)
	})

	// temporary commented out
	// ----------------------------------------
	//r.Post("/update", func(c *fiber.Ctx) error {
	//	return balances.UpdateBalance(c)
	//})

}
