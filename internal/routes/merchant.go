package routes

import (
	"github.com/dw-account-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initMerchantRoutes(router fiber.Router) {
	r := router.Group("/merchant")
	h := handlers.MerchantHandler{}

	r.Post("/register", func(c *fiber.Ctx) error {
		return h.Register(c)
	})

	r.Post("/unregister", func(c *fiber.Ctx) error {
		return h.Unregister(c)
	})

	r.Post("/all", func(c *fiber.Ctx) error {
		return h.GetMerchants(c)
	})

	r.Get("/", func(c *fiber.Ctx) error {
		return h.GetMerchantByID(c)
	})

	//r.Get("/uid/:uid", func(c *fiber.Ctx) error {
	//	return h.GetRegisteredAccountByUID(c)
	//})

	// ------------ TRX ------------

	r.Post("/balance/topup", func(c *fiber.Ctx) error {
		return h.BalanceTopup(c)
	})

	r.Post("/balance/inquiry", func(c *fiber.Ctx) error {
		return h.BalanceInquiry(c)
	})
}
