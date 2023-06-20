package handlers

import (
	"fmt"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/pkg/payload/request"
	"github.com/dw-account-service/pkg/tools"
	"github.com/dw-account-service/pkg/xlogger"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type MerchantHandler struct {
}

// isExists
func (m *MerchantHandler) isExists(account *entity.AccountBalance) bool {
	_, _, err := repository.AccountRepository.FindByMerchantID(account)
	if err != nil {
		// no document found, its mean it can be registered
		if err == mongo.ErrNoDocuments {
			return false
		}

		xlogger.Log.Println(err.Error())
		return true
	}
	return true
}

// isRegistered is a private function that check whether account id has been registered based on phoneNumber
func (m *MerchantHandler) isRegistered(account *entity.AccountBalance) bool {
	_, _, err := repository.AccountRepository.FindByMerchantStatus(account, true)
	if err != nil {
		// no document found, its mean it can be registered
		if err == mongo.ErrNoDocuments {
			return false
		}

		xlogger.Log.Println(err.Error())
		return true
	}
	return true
}

// isUnregistered is a private function that check whether account id has been unregistered or not
func (m *MerchantHandler) isUnregistered(ua *entity.UnregisterAccount) bool {

	acc := new(entity.AccountBalance)
	acc.MerchantID = ua.MerchantID
	acc.PartnerID = ua.PartnerID

	_, account, err := repository.AccountRepository.FindByMerchantStatus(acc, false)
	if err != nil {
		// no document found, its mean it can be unregistered
		if err == mongo.ErrNoDocuments {
			return false
		}

		xlogger.Log.Println(err.Error())
		return true
	}

	if account.(*entity.AccountBalance).Active == true {
		return false
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
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if m.isExists(payload) {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: "merchantId already exists, its probably in deactivated status. ",
			Data:    fiber.Map{"merchantId": payload.MerchantID},
		})
	}

	// validate request
	validation, err := tools.ValidateRequest(payload)
	if err != nil {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data: map[string]interface{}{
				"errors": validation,
			},
		})
	}

	// set default value for accountBalance document
	key, _ := tools.GenerateSecretKey()
	encryptedBalance, _ := tools.Encrypt([]byte(key), fmt.Sprintf("%016s", "0"))

	payload.SecretKey = key
	payload.Active = true
	payload.LastBalance = encryptedBalance
	payload.LastBalanceNumeric = 0
	payload.CreatedAt = time.Now().UnixMilli()
	payload.UpdatedAt = payload.CreatedAt

	code, id, err := repository.AccountRepository.InsertDocument(payload)
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	_, createdAccount, _ := repository.AccountRepository.FindByID(id, true)

	return c.Status(code).JSON(Responses{
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
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// check if already been unregistered
	if m.isUnregistered(u) {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: "account has already been deactivated",
			Data:    nil,
		})
	}

	code, err := repository.AccountRepository.DeactivateMerchant(u)
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// insert into accountDeactivated collection
	u.CreatedAt = time.Now().Format("2006-06-02 15:04:05")
	u.UpdatedAt = u.CreatedAt
	code, _, err = repository.AccountRepository.InsertDeactivatedAccount(u)

	ab := new(entity.AccountBalance)
	ab.PartnerID = u.PartnerID
	ab.MerchantID = u.MerchantID

	_, updatedAccount, _ := repository.AccountRepository.FindByMerchantStatus(ab, false)

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: "merchant successfully deactivated",
		Data:    updatedAccount,
	})
}

func (m *MerchantHandler) GetMerchants(c *fiber.Ctx) error {

	var req = new(request.PaginatedAccountRequest)

	// parse body payload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(Responses{
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
			return c.Status(400).JSON(Responses{
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

	code, merchants, total, pages, err := repository.AccountRepository.FindAllMerchant(req)
	if err != nil {
		return c.Status(code).JSON(ResponsePayloadPaginated{
			Success: false,
			Message: err.Error(),
			Data:    ResponsePayloadDataPaginated{},
		})
	}

	return c.Status(code).JSON(ResponsePayloadPaginated{
		Success: true,
		Message: msgResponse,
		Data: ResponsePayloadDataPaginated{
			Result:      merchants,
			Total:       total,
			PerPage:     req.Size,
			CurrentPage: req.Page,
			LastPage:    pages,
		},
	})
}

// GetMerchantByID is used to find registered account with active status = true
func (m *MerchantHandler) GetMerchantByID(c *fiber.Ctx) error {
	pid := c.Query("pid")
	mid := c.Query("mid")

	if pid == "" {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: "partnerId cannot be empty",
			Data:    nil,
		})
	}

	if mid == "" {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: "merchantId cannot be empty",
			Data:    nil,
		})
	}

	payload := new(entity.AccountBalance)
	payload.MerchantID = mid
	payload.PartnerID = pid

	code, account, err := repository.AccountRepository.FindByMerchantStatus(payload, true)

	if err != nil {
		errMsg := err.Error()
		if err == mongo.ErrNoDocuments {
			errMsg = "merchants not found or its already been deactivated"
		}
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: errMsg,
			Data:    nil,
		})
	}

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: "merchant fetched successfully ",
		Data:    account,
	})

}

// -------------------------- Transactions ------------------------------

func (m *MerchantHandler) BalanceTopup(c *fiber.Ctx) error {
	// TODO: merchant last balance topup
	// partnerId, merchantId, amount, partnerTransDate, partnerRefNo

	return nil
}
