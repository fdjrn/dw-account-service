package handlers

import (
	"context"
	"fmt"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/pkg/payload/request"
	"github.com/dw-account-service/pkg/tools"
	"github.com/dw-account-service/pkg/xlogger"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type AccountHandler struct {
}

func (a *AccountHandler) isExists(account *entity.AccountBalance) bool {
	_, err := repository.Account.FindOne(account)
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

func (a *AccountHandler) Register(c *fiber.Ctx) error {
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

	payload.Type = repository.AccountTypeRegular

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

	if a.isExists(payload) {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: "account already exists, or its probably in deactivated status. ",
			Data: map[string]interface{}{
				"partnerId":  payload.PartnerID,
				"merchantId": payload.MerchantID,
				"terminalId": payload.TerminalID,
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

	code, insertedId, err := repository.Account.InsertDocument(payload)
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	_, createdAccount, _ := repository.Account.FindByID(insertedId, true)

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: "account successfully registered",
		Data:    createdAccount,
	})
}

func (a *AccountHandler) Unregister(c *fiber.Ctx) error {

	// new UnregisterAccount struct
	uac := new(entity.UnregisterAccount)

	// parse body payload
	if err := c.BodyParser(uac); err != nil {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	code, err := repository.Account.DeactivateAccount(uac)
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// insert into accountDeactivated collection
	auditLog := time.Now().Format("2006-01-02 15:04:05")
	uac.Type = repository.AccountTypeRegular
	uac.CreatedAt = auditLog
	uac.UpdatedAt = auditLog
	code, _, err = repository.Account.InsertDeactivatedAccount(uac)

	updAccount := new(entity.AccountBalance)
	updAccount.PartnerID = uac.PartnerID
	updAccount.MerchantID = uac.MerchantID
	updAccount.TerminalID = uac.TerminalID

	result, _ := repository.Account.FindOne(updAccount)

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: "account successfully unregistered",
		Data:    result,
	})
}

func (a *AccountHandler) GetAllRegisteredAccountPaginated(c *fiber.Ctx) error {

	var req = new(request.PaginatedAccountRequest)

	// parse body payload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	msgResponse := "accounts successfully fetched"
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
		msgResponse = fmt.Sprintf("%s account successfully fetched", req.Status)
	}

	// set default value
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}

	code, accounts, total, pages, err := repository.Account.FindAllPaginated(req)
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
			Result:      accounts,
			Total:       total,
			PerPage:     req.Size,
			CurrentPage: req.Page,
			LastPage:    pages,
		},
	})
}

func (a *AccountHandler) GetActiveAccountByID(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	code, account, err := repository.Account.FindByID(id, true)

	if err != nil {
		errMsg := err.Error()
		if err == mongo.ErrNoDocuments {
			errMsg = "account not found or its already been unregistered"
		}
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: errMsg,
			Data:    nil,
		})
	}

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: "account fetched successfully ",
		Data:    account,
	})

}

func (a *AccountHandler) GetActiveAccountByUniqueID(c *fiber.Ctx) error {
	code, account, err := repository.Account.FindByActiveStatus(c.Params("uid"), true)

	if err != nil {
		errMsg := err.Error()
		if err == mongo.ErrNoDocuments {
			errMsg = "account not found or its already been unregistered"
		}

		return c.Status(code).JSON(Responses{
			Success: false,
			Message: errMsg,
			Data:    nil,
		})
	}

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: "account fetched successfully ",
		Data:    account,
	})

}

func (a *AccountHandler) GetAccountDetail(c *fiber.Ctx) error {

	// new account struct
	payload := new(entity.AccountBalance)

	// parse body payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate request
	validation, err := tools.ValidateAccountDetailRequest(payload)
	if err != nil {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data: map[string]interface{}{
				"errors": validation,
			},
		})
	}

	account, err := repository.Account.FindOne(payload)

	if err != nil {
		errMsg := err.Error()
		if err == mongo.ErrNoDocuments {
			errMsg = "account not found or its already been unregistered"
		}

		return c.Status(500).JSON(Responses{
			Success: false,
			Message: errMsg,
			Data:    nil,
		})
	}

	return c.Status(200).JSON(Responses{
		Success: true,
		Message: "account fetched successfully ",
		Data:    account,
	})

}

// ----------------- utilities -----------------

func (a *AccountHandler) UpdateMerchantAndTerminalForAccount(c *fiber.Ctx) error {

	filter := bson.D{
		{"merchantId", primitive.Null{}},
		{"terminalId", primitive.Null{}},
		{"type", 1},
		{"uniqueId", bson.D{{"$ne", primitive.Null{}}}},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	cursor, err := db.Mongo.Collection.Account.Find(
		ctx,
		filter,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": true,
			"message": err.Error(),
			"count":   0,
			"data":    nil,
		})
	}

	var accounts []entity.AccountBalance
	if err = cursor.All(context.TODO(), &accounts); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": true,
			"message": err.Error(),
			"count":   0,
			"data":    nil,
		})
	}

	//var arrAccount []entity.AccountBalance
	var successCount int64
	for _, account := range accounts {
		str := strings.Split(account.UniqueID, "_")
		account.TerminalID = str[0]
		account.MerchantID = str[1]

		result, err2 := db.Mongo.Collection.Account.UpdateOne(context.TODO(),
			filter,
			bson.D{
				{"$set", bson.D{
					{"terminalId", str[0]},
					{"merchantId", str[1]},
				}},
			})
		if err2 != nil {
			xlogger.Log.Println("err: ", err2.Error())
		}

		//arrAccount = append(arrAccount, account)

		successCount = successCount + result.ModifiedCount

	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "ok",
		"count":   fmt.Sprintf("%d account has been successfully updated", successCount),
		"data":    nil,
	})

}

func (a *AccountHandler) SyncBalance(c *fiber.Ctx) error {

	filter := bson.D{
		//{"merchantId", nil},
		//{"terminalId", nil},
		//{"uniqueId", bson.D{{"$ne", nil}}},
		{"lastBalanceNumeric", primitive.Null{}},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	cursor, err := db.Mongo.Collection.Account.Find(
		ctx,
		filter,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": true,
			"message": err.Error(),
			"count":   0,
			"data":    nil,
		})
	}

	var accounts []entity.AccountBalance
	if err = cursor.All(context.TODO(), &accounts); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": true,
			"message": err.Error(),
			"count":   0,
			"data":    nil,
		})
	}

	//var arrAccount []entity.AccountBalance
	var successCount int64
	for _, account := range accounts {
		currentBalance, _ := tools.DecryptAndConvert([]byte(account.SecretKey), account.LastBalance)

		//str := strings.Split(account.UniqueID, "_")
		//account.TerminalID = str[0]
		//account.MerchantID = str[1]

		result, err2 := db.Mongo.Collection.Account.UpdateOne(context.TODO(),
			filter,
			bson.D{
				{"$set", bson.D{
					{"lastBalanceNumeric", int64(currentBalance)},
				}},
			})
		if err2 != nil {
			xlogger.Log.Println("err: ", err2.Error())
		}

		//	arrAccount = append(arrAccount, account)

		successCount = successCount + result.ModifiedCount

	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "ok",
		"count":   fmt.Sprintf("%d account has been successfully updated", successCount),
		"data":    nil,
	})
}
