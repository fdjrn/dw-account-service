package handlers

import (
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/internal/utilities"
	"github.com/gofiber/fiber/v2"
)

type BalanceHandler struct {
	repo repository.BalanceRepository
}

func NewBalanceHandler() BalanceHandler {
	return BalanceHandler{repo: repository.NewBalanceRepository()}
}

func (b *BalanceHandler) Inquiry(c *fiber.Ctx, isMerchant bool) error {
	payload := new(entity.InquiryBalance)
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	payload.Type = utilities.AccountTypeMerchant
	if !isMerchant {
		if payload.TerminalID == "" {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "terminalId cannot be empty",
				Data:    nil,
			})
		}

		payload.Type = utilities.AccountTypeRegular
	}

	// set empty terminalId if its merchant
	if payload.Type == utilities.AccountTypeMerchant {
		payload.TerminalID = ""
	}

	b.repo.Entity = payload

	err := b.repo.GetLastBalance()
	if err != nil {
		return SendDefaultErrResponse("failed to inquiry last balance on current merchant, ", err, c)
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "balance successfully fetched",
		Data:    b.repo.Entity,
	})
}

func (b *BalanceHandler) MerchantBalanceSummary(c *fiber.Ctx) error {
	return nil
}
