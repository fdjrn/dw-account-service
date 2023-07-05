package routes

import (
	"github.com/dw-account-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initBalanceRoutes(router fiber.Router) {
	balanceHandler := handlers.NewBalanceHandler()

	r := router.Group("/account")
	r.Post("/balance/inquiry", func(c *fiber.Ctx) error {
		return balanceHandler.Inquiry(c, false)
	})

	// ---------------------------------------------------------------

	r2 := router.Group("/merchant")
	r2.Post("/balance/inquiry", func(c *fiber.Ctx) error {
		return balanceHandler.Inquiry(c, true)
	})

	r2.Post("/balance/summary", func(c *fiber.Ctx) error {
		return balanceHandler.MerchantBalanceSummary(c)
	})

	// balance transaction
	// ---------------------------------------------------------------
	//r.Post("/topup", func(c *fiber.Ctx) error {
	//	return balanceHandler.TopupBalance(c)
	//})

	//r.Post("/topup-merchant", func(c *fiber.Ctx) error {
	//	return h.TopupMerchantBalance(c)
	//})

	//r.Post("/deduct", func(c *fiber.Ctx) error {
	//	return balanceHandler.DeductBalance(c)
	//})

	// temporary commented out
	// ----------------------------------------
	//r.Post("/update", func(c *fiber.Ctx) error {
	//	return balances.UpdateBalance(c)
	//})

}
