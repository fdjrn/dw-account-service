package routes

import (
	"github.com/dw-account-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initBalanceRoutes(router fiber.Router) {
	r := router.Group("/account/balance")
	h := handlers.BalanceHandler{}

	r.Get("/inquiry/:uid", func(c *fiber.Ctx) error {
		return h.InquiryBalance(c)
	})

	// balance transaction
	// ---------------------------------------------------------------
	r.Post("/topup", func(c *fiber.Ctx) error {
		return h.TopupBalance(c)
	})

	//r.Post("/topup-merchant", func(c *fiber.Ctx) error {
	//	return h.TopupMerchantBalance(c)
	//})

	r.Post("/deduct", func(c *fiber.Ctx) error {
		return h.DeductBalance(c)
	})

	// temporary commented out
	// ----------------------------------------
	//r.Post("/update", func(c *fiber.Ctx) error {
	//	return balances.UpdateBalance(c)
	//})

}
