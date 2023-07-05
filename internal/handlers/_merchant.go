package handlers

import (
	"fmt"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/internal/handlers/validator"
	"github.com/dw-account-service/internal/utilities"
	"github.com/dw-account-service/internal/utilities/crypt"
	"github.com/dw-account-service/internal/utilities/str"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"time"
)

type MerchantHandler struct {
}

// isExists
func (m *MerchantHandler) isExists(account *entity.AccountBalance) bool {
	_, err := repository.Account.FindOne()
	if err != nil {
		// no document found, its mean it can be registered
		if err == mongo.ErrNoDocuments {
			return false
		}

		utilities.Log.Println(err.Error())
		return true
	}
	return true
}

// Register is a function that used to insert new document into collection and set active status to true.
func (m *MerchantHandler) Register(c *fiber.Ctx) error {
	var err error

	// new account struct
	payload := new(entity.AccountBalance)

	// parse body payload
	if err = c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	payload.Type = repository.AccountTypeMerchant

	if m.isExists(payload) {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: "merchantId already exists, or its probably in deactivated status. ",
			Data:    fiber.Map{"merchantId": payload.MerchantID, "partnerId": payload.PartnerID},
		})
	}

	// validate request
	validation, err := validator.ValidateRequest(payload)
	if err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data: map[string]interface{}{
				"errors": validation,
			},
		})
	}

	// set default value for accountBalance document
	key, _ := crypt.GenerateSecretKey()
	encryptedBalance, _ := crypt.Encrypt([]byte(key), fmt.Sprintf("%016s", "0"))

	payload.SecretKey = key
	payload.Active = true

	payload.LastBalance = encryptedBalance
	payload.LastBalanceNumeric = 0
	payload.CreatedAt = time.Now().UnixMilli()
	payload.UpdatedAt = payload.CreatedAt

	id, err := repository.Account.Create(payload)
	if err != nil {
		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	_, createdAccount, _ := repository.Account.FindByID(id, true)

	return c.Status(code).JSON(entity.Responses{
		Success: true,
		Message: "merchant successfully registered",
		Data:    createdAccount,
	})
}

// Unregister is a function that used to change active status to false (unregistered)
func (m *MerchantHandler) Unregister(c *fiber.Ctx) error {

	// new u struct
	u := new(entity.UnregisterAccount)

	// parse body payload
	if err := c.BodyParser(u); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	code, err := repository.Account.DeactivateMerchant(u)
	if err != nil {
		return c.Status(code).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// insert into accountDeactivated collection
	auditLog := time.Now().Format("2006-01-02 15:04:05")
	u.Type = repository.AccountTypeMerchant
	u.CreatedAt = auditLog
	u.UpdatedAt = auditLog
	code, _, err = repository.Account.InsertDeactivatedAccount(u)

	updAccount := new(entity.AccountBalance)
	updAccount.PartnerID = u.PartnerID
	updAccount.MerchantID = u.MerchantID

	//_, updatedAccount, _ := repository.Account.FindByMerchantStatus(ab, false)
	result, _ := repository.Account.FindOne(updAccount)

	return c.Status(code).JSON(entity.Responses{
		Success: true,
		Message: "merchant successfully deactivated",
		Data:    result,
	})
}

func (m *MerchantHandler) GetMerchants(c *fiber.Ctx) error {

	var req = new(entity.PaginatedAccountRequest)

	// parse body payload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	msgResponse := "merchants successfully fetched"
	req.Status = strings.ToLower(req.Status)

	validStatus := map[string]interface{}{"all": 0, "active": 1, "deactivated": 2}
	if req.Status != "" {
		if _, ok := validStatus[req.Status]; !ok {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "invalid status value. its only accept all, active or deactivated",
				Data:    nil,
			})
		}
		msgResponse = fmt.Sprintf("%s merchants successfully fetched", req.Status)
	}

	// set default value
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}

	code, merchants, total, pages, err := repository.Account.FindAllMerchant(req)
	if err != nil {
		return c.Status(code).JSON(entity.ResponsePayloadPaginated{
			Success: false,
			Message: err.Error(),
			Data:    entity.ResponsePayloadDataPaginated{},
		})
	}

	return c.Status(code).JSON(entity.ResponsePayloadPaginated{
		Success: true,
		Message: msgResponse,
		Data: entity.ResponsePayloadDataPaginated{
			Result:      merchants,
			Total:       total,
			PerPage:     req.Size,
			CurrentPage: req.Page,
			LastPage:    pages,
		},
	})
}

// GetMerchantDetail is used to find registered merchant with active status = true
func (m *MerchantHandler) GetMerchantDetail(c *fiber.Ctx) error {

	// new account struct
	payload := new(entity.AccountBalance)
	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	payload.Type = repository.AccountTypeMerchant
	payload.Active = true

	account, err := repository.Account.FindOne(payload)

	if err != nil {
		errMsg := err.Error()
		if err == mongo.ErrNoDocuments {
			errMsg = "merchants not found or its already been deactivated"
		}
		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: errMsg,
			Data:    nil,
		})
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "merchant fetched successfully ",
		Data:    account,
	})

}



// -------------------------- Transactions ------------------------------

func (m *MerchantHandler) BalanceTopup(c *fiber.Ctx) error {

	// new MerchantTrxRequest struct
	payload := new(entity.MerchantTrxRequest)

	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate request
	msg, err := validator.ValidateRequest(payload)
	if err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data: map[string]interface{}{
				"errors": msg,
			},
		})
	}

	// 1. check used partnerRefNumber
	if repository.Topup.IsUsedPartnerRefNumber(payload.PartnerRefNumber) {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: "partnerRefNumber already used",
			Data: fiber.Map{
				"partnerRefNumber": payload.PartnerRefNumber,
			},
		})
	}

	// 2. inquiry balance
	code, balance, err := repository.Balance.MerchantInquiry(entity.BalanceInquiry{
		MerchantID: payload.MerchantID, PartnerID: payload.PartnerID,
	})
	if err != nil {
		return c.Status(code).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	tBalance := new(entity.BalanceTopUp)

	// 3. add last balance with amount of topup
	//tBalance.UniqueID = "-"
	tBalance.PartnerID = payload.PartnerID
	tBalance.MerchantID = payload.MerchantID
	//tBalance.TerminalID = "-"
	//tBalance.VoucherCode = "-"
	//tBalance.VoucherAmount = payload.Amount
	tBalance.Amount = payload.Amount
	tBalance.PartnerRefNumber = payload.PartnerRefNumber
	tBalance.PartnerTransDate = payload.PartnerTransDate
	tBalance.TransDate = time.Now().UnixMilli()
	tBalance.TransNumber = str.GenerateTransNumber()
	tBalance.ReceiptNumber = str.GenerateReceiptNumber(utilities.TransTopUp, "")

	// 4. encrypt result addition
	tBalance.LastBalance = balance.CurrentBalance + int64(payload.Amount)
	strLastBalance := strconv.FormatInt(tBalance.LastBalance, 10)
	tBalance.LastBalanceEncrypted, _ = crypt.Encrypt([]byte(balance.SecretKey), fmt.Sprintf("%016s", strLastBalance))

	tBalance.Status = repository.TrxStatusSuccess
	tBalance.CreatedAt = tBalance.TransDate
	tBalance.UpdatedAt = tBalance.TransDate

	// 5. insert topup document
	code, err = repository.Topup.CreateTopupDocument(tBalance)
	if err != nil {
		return c.Status(code).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// 6. do update document on AccountRepository Collection
	code, err = repository.Balance.UpdateMerchantBalance(tBalance)
	if err != nil {
		return c.Status(code).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// TODO: 7. publish transaction to broker
	//payload, err := json.Marshal(t)
	//if err != nil {
	//	log.Println("cannot marshal payload: ", err.Error())
	//}
	//
	//_ = kafka.ProduceMsg(kafka.TopUpTopic, payload)

	return c.Status(code).JSON(entity.Responses{
		Success: true,
		Message: "account balance has been top-up successfully",
		Data:    tBalance,
	})
}

func (m *MerchantHandler) BalanceInquiry(c *fiber.Ctx) error {

	var payload entity.BalanceInquiry

	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	code, account, err := repository.Balance.MerchantInquiry(payload)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(code).JSON(entity.Responses{
				Success: false,
				Message: "balance not found or it has been unregistered",
				Data:    nil,
			})
		}

		return c.Status(code).JSON(entity.Responses{
			Success: false,
			Message: "failed to inquiry last balance on current merchant",
			Data:    nil,
		})
	}

	return c.Status(code).JSON(entity.Responses{
		Success: true,
		Message: "balance successfully fetched",
		Data:    account,
	})
}

func (m *MerchantHandler) BalanceDeduct(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "under construction",
		"data":    nil,
	})
}

func (m *MerchantHandler) BalanceDistribution(c *fiber.Ctx) error {

	// new MerchantTrxRequest struct
	payload := new(entity.MerchantTrxRequest)

	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate request
	msg, err := validator.ValidateRequest(payload)
	if err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data: map[string]interface{}{
				"errors": msg,
			},
		})
	}

	// TODO 1: check used partnerRefNumber

	// 2: inquiry balance
	code, balance, err := repository.Balance.MerchantInquiry(entity.BalanceInquiry{
		MerchantID: payload.MerchantID, PartnerID: payload.PartnerID,
	})
	if err != nil {
		return c.Status(code).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// 3: check total member in merchant
	countDoc, err := repository.Account.CountMerchantMember(
		entity.AccountBalance{
			PartnerID:  payload.PartnerID,
			MerchantID: payload.MerchantID,
			Active:     true,
		})

	if balance.LastBalanceNumeric < (payload.Amount * countDoc) {
		return c.Status(500).JSON(entity.Responses{
			Success: false,
			Message: "insufficient amount to distribute balance to others",
			Data:    nil,
		})
	}

	// TODO 4: update merchant balance amount

	// TODO 5: publish to kafka for history and summary

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "under construction",
		"data":    nil,
	})
}
