package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/internal/kafka"
	"github.com/dw-account-service/pkg/tools"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"time"
)

func InquiryBalance(c *fiber.Ctx) error {
	code, account, err := repository.BalanceRepository.Inquiry(c.Params("uid"))
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: "balance successfully fetched",
		Data:    account,
	})
}

func TopupBalance(c *fiber.Ctx) error {
	// new BalanceTopup struct
	t := new(entity.BalanceTopUp)

	// parse body payload
	if err := c.BodyParser(t); err != nil {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// 1. check used exRefNumber
	if repository.TopupRepository.IsUsedExRefNumber(t.ExRefNumber) {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: "exRefNumber already used",
			Data:    t,
		})
	}

	// 2. inquiry balance
	code, b, err := repository.BalanceRepository.Inquiry(t.MDLUniqueID)
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// 3. add last balance with amount of topup
	t.LastBalance = b.CurrentBalance + int64(t.Amount)

	// 4. encrypt result addition
	strLastBalance := strconv.FormatInt(t.LastBalance, 10)
	t.LastBalanceEncrypted, _ = tools.Encrypt([]byte(b.SecretKey), fmt.Sprintf("%016s", strLastBalance))

	t.ReceiptNumber = tools.GenerateReceiptNumber(tools.TransTopUp, "")
	t.CreatedAt = time.Now().UnixMilli()
	t.UpdatedAt = t.CreatedAt

	// 5. insert topup document
	code, err = repository.TopupRepository.CreateTopupDocument(t)
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// 6. do update document on Account Collection
	code, err = repository.BalanceRepository.UpdateBalance(t.MDLUniqueID, t.LastBalanceEncrypted)
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// 7. publish transaction to broker
	payload, err := json.Marshal(t)
	if err != nil {
		log.Println("cannot marshal payload: ", err.Error())
	}

	_ = kafka.ProduceMsg(kafka.TopUpTopic, payload)

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: "account balance has been top-up successfully",
		Data:    t,
	})
}

func DeductBalance(c *fiber.Ctx) error {
	// new BalanceDeduction struct
	d := new(entity.BalanceDeduction)

	// parse body payload
	if err := c.BodyParser(d); err != nil {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// 1. Inquiry Balance
	_, b, err := repository.BalanceRepository.Inquiry(d.MDLUniqueID)
	if err != nil {
		return c.Status(500).JSON(Responses{
			Success: false,
			Message: "failed to inquiry last balance on current account",
			Data:    nil,
		})
	}
	// 2. Jika saldo cukup, maka lanjutkan proses pengurangan saldo (pembayaran)
	if b.CurrentBalance < int64(d.Amount) {
		return c.Status(500).JSON(Responses{
			Success: false,
			Message: "current balance is less than current transaction amount",
			Data:    &b,
		})
	}

	d.LastBalance = b.CurrentBalance - int64(d.Amount)

	// 3. Encrypt hasil pengurangan
	strLastBalance := strconv.FormatInt(d.LastBalance, 10)
	d.LastBalanceEncrypted, _ = tools.Encrypt([]byte(b.SecretKey), fmt.Sprintf("%016s", strLastBalance))

	// 4. Update document
	code, err := repository.BalanceRepository.UpdateBalance(d.MDLUniqueID, d.LastBalanceEncrypted)
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	// 5. Fetch updated document
	_, b, _ = repository.BalanceRepository.Inquiry(d.MDLUniqueID)

	// 6. Populate response data
	// untuk PoC masih Hardcoded dulu...

	// TransType --> 1: Pembelian Konten E-Course | 2: TBD
	// d.TransType = 1
	// d.Description = "Pembelian Konten e-Course MyDigiLearn"
	d.ReceiptNumber = tools.GenerateReceiptNumber(tools.TransPayment, "")
	d.LastBalance = b.CurrentBalance

	// 7. Publish payment/deduction message to broker
	payload, err := json.Marshal(d)
	if err != nil {
		log.Println("cannot marshal payload: ", err.Error())
	}

	_ = kafka.ProduceMsg(kafka.DeductTopic, payload)

	return c.Status(200).JSON(Responses{
		Success: true,
		Message: "balance has been successfully deducted",
		Data:    d,
	})
}
