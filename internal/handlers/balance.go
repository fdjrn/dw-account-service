package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
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

	payload.Type = repository.AccountTypeMerchant
	if !isMerchant {
		if payload.TerminalID == "" {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "terminalId cannot be empty",
				Data:    nil,
			})
		}

		payload.Type = repository.AccountTypeRegular
	}

	// set empty terminalId if its merchant
	if payload.Type == repository.AccountTypeMerchant {
		payload.TerminalID = ""
	}

	b.repo.Entity = payload

	err := b.repo.GetLastBalance()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(entity.Responses{
				Success: false,
				Message: "balance not found or it has been unregistered",
				Data:    nil,
			})
		}

		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: "failed to inquiry last balance on current merchant",
			Data:    nil,
		})
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "balance successfully fetched",
		Data:    b.repo.Entity,
	})
}

// MerchantBalanceSummary is used to get balance summary based on period
func (b *BalanceHandler) MerchantBalanceSummary(c *fiber.Ctx) error {
	return nil
}

//func (b *BalanceHandler) TopupBalance(c *fiber.Ctx) error {
//	// new BalanceTopup struct
//	t := new(entity.BalanceTopUp)
//
//	// parse body payload
//	if err := c.BodyParser(t); err != nil {
//		return c.Status(400).JSON(entity.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data:    nil,
//		})
//	}
//
//	// validate request
//	m, err := validator.ValidateRequest(t)
//	if err != nil {
//		return c.Status(400).JSON(entity.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data: map[string]interface{}{
//				"errors": m,
//			},
//		})
//	}
//
//	// 1. check used exRefNumber
//	if repository.Topup.IsUsedExRefNumber(t.ExRefNumber) {
//		return c.Status(400).JSON(entity.Responses{
//			Success: false,
//			Message: "exRefNumber already used",
//			Data:    t,
//		})
//	}
//
//	// 2. inquiry balance
//	code, balance, err := repository.Balance.Inquiry(t.UniqueID)
//	if err != nil {
//		return c.Status(code).JSON(entity.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data:    nil,
//		})
//	}
//
//	// 3. add last balance with amount of topup
//	t.LastBalance = balance.CurrentBalance + int64(t.Amount)
//
//	// 4. encrypt result addition
//	strLastBalance := strconv.FormatInt(t.LastBalance, 10)
//	t.LastBalanceEncrypted, _ = crypt.Encrypt([]byte(balance.SecretKey), fmt.Sprintf("%016s", strLastBalance))
//
//	t.ReceiptNumber = str.GenerateReceiptNumber(utilities.TransTopUp, "")
//	t.CreatedAt = time.Now().UnixMilli()
//	t.UpdatedAt = t.CreatedAt
//
//	// 5. insert topup document
//	code, err = repository.Topup.CreateTopupDocument(t)
//	if err != nil {
//		return c.Status(code).JSON(entity.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data:    nil,
//		})
//	}
//
//	// 6. do update document on AccountRepository Collection
//	code, err = repository.Balance.UpdateBalance(t.UniqueID, t.LastBalanceEncrypted)
//	if err != nil {
//		return c.Status(code).JSON(entity.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data:    nil,
//		})
//	}
//
//	// 7. publish transaction to broker
//	//payload, err := json.Marshal(t)
//	//if err != nil {
//	//	log.Println("cannot marshal payload: ", err.Error())
//	//}
//
//	//_ = kafka.ProduceMsg(kafka.TopUpTopic, payload)
//
//	return c.Status(code).JSON(entity.Responses{
//		Success: true,
//		Message: "account balance has been top-up successfully",
//		Data:    t,
//	})
//}

//func (b *BalanceHandler) DeductBalance(c *fiber.Ctx) error {
//	// new BalanceDeduction struct
//	d := new(entity.BalanceDeduction)
//
//	// parse body payload
//	if err := c.BodyParser(d); err != nil {
//		return c.Status(400).JSON(entity.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data:    nil,
//		})
//	}
//
//	// validate request
//	m, err := validator.ValidateRequest(d)
//	if err != nil {
//		return c.Status(400).JSON(entity.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data: map[string]interface{}{
//				"errors": m,
//			},
//		})
//	}
//
//	// 1. Inquiry BalanceRepository
//	code, balance, err := repository.Balance.Inquiry(d.UniqueID)
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return c.Status(code).JSON(entity.Responses{
//				Success: false,
//				Message: "balance not found or it has been unregistered",
//				Data:    nil,
//			})
//		}
//
//		return c.Status(500).JSON(entity.Responses{
//			Success: false,
//			Message: "failed to inquiry last balance on current account",
//			Data:    nil,
//		})
//	}
//
//	// TODO simpan current last balance, digunakan untuk proses rollback transaction
//	// jika terjadi kegagalan produce message
//
//	// 2. Jika saldo cukup, maka lanjutkan proses pengurangan saldo (pembayaran)
//	if balance.CurrentBalance < int64(d.Amount) {
//		return c.Status(500).JSON(entity.Responses{
//			Success: false,
//			Message: "current balance is less than current transaction amount",
//			Data:    &b,
//		})
//	}
//
//	d.LastBalance = balance.CurrentBalance - int64(d.Amount)
//
//	// 3. Encrypt hasil pengurangan
//	strLastBalance := strconv.FormatInt(d.LastBalance, 10)
//	d.LastBalanceEncrypted, _ = crypt.Encrypt([]byte(balance.SecretKey), fmt.Sprintf("%016s", strLastBalance))
//
//	// 4. Update document
//	code, err = repository.Balance.UpdateBalance(d.UniqueID, d.LastBalanceEncrypted)
//	if err != nil {
//		return c.Status(code).JSON(entity.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data:    nil,
//		})
//	}
//
//	// 5. Fetch updated document
//	_, balance, _ = repository.Balance.Inquiry(d.UniqueID)
//
//	// 6. set default value
//	// untuk PoC masih Hardcoded dulu...
//	// TransType --> 1: Pembelian Konten E-Course | 2: TBD
//	d.ReceiptNumber = str.GenerateReceiptNumber(utilities.TransPayment, "")
//	d.LastBalance = balance.CurrentBalance
//
//	// 7. Publish payment/deduction message to broker
//	//payload, err := json.Marshal(d)
//	//if err != nil {
//	//	log.Println("cannot marshal payload: ", err.Error())
//	//}
//	//
//	//_ = kafka.ProduceMsg(kafka.DeductTopic, payload)
//
//	return c.Status(200).JSON(entity.Responses{
//		Success: true,
//		Message: "balance has been successfully deducted",
//		Data:    d,
//	})
//}

//func (b *BalanceHandler) BalanceInquiry(c *fiber.Ctx) error {
//
//	uid := c.Params("uid")
//	if uid == "" {
//		return c.Status(400).JSON(entity.Responses{
//			Success: false,
//			Message: "uniqueId cannot be empty",
//			Data:    nil,
//		})
//	}
//
//	code, account, err := repository.Balance.InquiryBalance(uid)
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return c.Status(code).JSON(entity.Responses{
//				Success: false,
//				Message: "balance not found or it has been unregistered",
//				Data:    nil,
//			})
//		}
//
//		return c.Status(code).JSON(entity.Responses{
//			Success: false,
//			Message: "failed to inquiry last balance on current account",
//			Data:    nil,
//		})
//	}
//
//	return c.Status(code).JSON(entity.Responses{
//		Success: true,
//		Message: "balance successfully fetched",
//		Data:    account,
//	})
//}

//func (b *BalanceHandler) MerchantBalanceInquiry(c *fiber.Ctx) error {
//
//	var payload entity.BalanceInquiry
//
//	// parse body payload
//	if err := c.BodyParser(&payload); err != nil {
//		return c.Status(400).JSON(entity.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data:    nil,
//		})
//	}
//
//	code, account, err := repository.Balance.MerchantInquiryBalance(payload)
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return c.Status(code).JSON(entity.Responses{
//				Success: false,
//				Message: "balance not found or it has been unregistered",
//				Data:    nil,
//			})
//		}
//
//		return c.Status(code).JSON(entity.Responses{
//			Success: false,
//			Message: "failed to inquiry last balance on current merchant",
//			Data:    nil,
//		})
//	}
//
//	return c.Status(code).JSON(entity.Responses{
//		Success: true,
//		Message: "balance successfully fetched",
//		Data:    account,
//	})
//}
