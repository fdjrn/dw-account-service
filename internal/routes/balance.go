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

	//r2.Post("/balance/summary", func(c *fiber.Ctx) error {
	//	return balanceHandler.MerchantBalanceSummary(c)
	//})

}
