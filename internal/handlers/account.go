package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/internal/handlers/validator"
	"github.com/dw-account-service/internal/utilities"
	"github.com/dw-account-service/internal/utilities/crypt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type AccountHandler struct {
	repo        repository.AccountRepository
	balanceRepo repository.BalanceRepository
}

func NewAccountHandler() AccountHandler {
	return AccountHandler{
		repo:        repository.NewAccountRepository(),
		balanceRepo: repository.NewBalanceRepository(),
	}
}

func (a *AccountHandler) existsAccount() (bool, error) {
	_, err := a.repo.FindOne()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a *AccountHandler) Register(c *fiber.Ctx) error {
	var err error

	// new account struct
	payload := new(entity.AccountBalance)

	// parse body payload
	if err = c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
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
	payload.UniqueID = fmt.Sprintf("%s%s", payload.MerchantID, payload.TerminalID)
	if payload.Type == utilities.AccountTypeMerchant {
		payload.UniqueID = ""
		payload.TerminalID = ""
		payload.TerminalName = ""
	}

	payload.LastBalance = encryptedBalance
	payload.LastBalanceNumeric = 0
	payload.CreatedAt = time.Now().UnixMilli()
	payload.UpdatedAt = payload.CreatedAt

	a.repo.Entity = payload

	exists, err := a.existsAccount()
	if !exists && err != nil {
		return SendDefaultErrResponse("failed to validate existing account, ", err, c)
	}

	if exists {
		responseData := map[string]interface{}{"partnerId": payload.PartnerID, "merchantId": payload.MerchantID}

		if payload.Type == utilities.AccountTypeRegular {
			responseData["terminalId"] = payload.TerminalID
		}

		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: "account already exists, or its probably in deactivated status. ",
			Data:    responseData,
		})
	}

	insertedId, err := a.repo.Create()
	if err != nil {
		return SendDefaultErrResponse("", err, c)
	}

	createdAccount, err := a.repo.FindByID(insertedId)
	if err != nil {
		return SendDefaultErrResponse("cannot fetch current registered account, ", err, c)
	}

	return c.Status(201).JSON(entity.Responses{
		Success: true,
		Message: "registration successful",
		Data:    createdAccount,
	})
}

func (a *AccountHandler) Unregister(c *fiber.Ctx) error {

	// new UnregisterAccount struct
	payload := new(entity.UnregisterAccount)

	// parse body payload
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// check is valid account
	a.repo.Entity.PartnerID = payload.PartnerID
	a.repo.Entity.MerchantID = payload.MerchantID
	a.repo.Entity.TerminalID = payload.TerminalID

	exists, err := a.existsAccount()
	if !exists {
		if err != nil {
			return SendDefaultErrResponse("failed to validate existing account, ", err, c)
		}

		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: "account not found",
			Data:    nil,
		})
	}

	err = a.repo.DeactivateAccount(payload)
	if err != nil {
		return SendDefaultErrResponse("", err, c)
	}

	// insert into accountDeactivated collection
	auditLog := time.Now().Format("2006-01-02 15:04:05")
	payload.CreatedAt = auditLog
	payload.UpdatedAt = auditLog

	//var acc entity.AccountBalance
	doc, _ := a.repo.FindOne()
	payload.Type = doc.Type
	payload.UniqueID = doc.UniqueID

	_, err = a.repo.InsertDeactivatedAccount(payload)
	if err != nil {
		return SendDefaultErrResponse("failed on insert deactivated account data, ", err, c)
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "deactivation successful",
		Data:    doc,
	})
}

func (a *AccountHandler) GetAccountByID(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	account, err := a.repo.FindByID(id)

	if err != nil {
		return SendDefaultErrResponse("failed to fetch account, ", err, c)
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "account fetched successfully ",
		Data:    account,
	})

}

func (a *AccountHandler) GetAccount(c *fiber.Ctx) error {

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

	a.repo.Entity = payload
	account, err := a.repo.FindOne()

	if err != nil {
		return SendDefaultErrResponse("failed to fetch account, ", err, c)
	}

	return c.Status(200).JSON(entity.Responses{
		Success: true,
		Message: "account fetched successfully ",
		Data:    account,
	})

}

func (a *AccountHandler) GetAccountsPaginated(c *fiber.Ctx) error {

	var req = new(entity.PaginatedAccountRequest)

	// parse body payload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(entity.Responses{
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
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "invalid status value. its only accept all, active or deactivated",
				Data:    nil,
			})
		}
		msgResponse = fmt.Sprintf("%s account successfully fetched", req.Status)
	}

	// set default param value
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}

	accounts, total, pages, err := a.repo.FindAllPaginated(req)
	if err != nil {
		return SendDefaultPaginationErrResponse("", err, c)
	}

	return c.Status(200).JSON(entity.PaginatedResponse{
		Success: true,
		Message: msgResponse,
		Data: entity.PaginatedDetailResponse{
			Result: accounts,
			Total:  total,
			Pagination: entity.PaginationInfo{
				PerPage:     req.Size,
				CurrentPage: req.Page,
				LastPage:    pages,
			},
		},
	})
}

// ----------------- Merchants -----------------

func (a *AccountHandler) GetMerchantMembers(c *fiber.Ctx, isPeriod bool) error {
	var err error

	// new account struct
	payload := new(entity.PaginatedAccountRequest)
	// parse body payload
	if err = c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate periods parameter
	if isPeriod {
		payload.Periods.StartDate, err = time.ParseInLocation(
			"20060102150405",
			fmt.Sprintf("%s%s", payload.Periods.Start, "000000"),
			time.Now().Location(),
		)
		if err != nil {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "invalid start periods",
				Data:    nil,
			})
		}

		payload.Periods.EndDate, err = time.ParseInLocation(
			"20060102150405",
			fmt.Sprintf("%s%s", payload.Periods.End, "235959"),
			time.Now().Location(),
		)

		if err != nil {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "invalid end periods",
				Data:    nil,
			})
		}

		if payload.Periods.EndDate.Before(payload.Periods.StartDate) {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "end period cannot be less than start period",
				Data:    nil,
			})
		}
	}

	// set type for payload validation
	payload.Type = utilities.AccountTypeMerchant

	// validate request
	if payload.PartnerID == "" {
		return c.Status(400).JSON(entity.PaginatedResponseMembers{
			Success: false,
			Message: "partnerId cannot be empty",
			Data:    nil,
		})
	}

	if payload.MerchantID == "" {
		return c.Status(400).JSON(entity.PaginatedResponseMembers{
			Success: false,
			Message: "merchantId cannot be empty",
			Data:    nil,
		})
	}

	// set default value
	if payload.Page == 0 {
		payload.Page = 1
	}

	if payload.Size == 0 {
		payload.Size = 10
	}

	// re-apply type for filter condition
	payload.Type = utilities.AccountTypeRegular

	members, total, pages, err := a.repo.FindMembersPaginated(payload, isPeriod)

	if err != nil {
		return SendDefaultPaginationErrResponse("cannot fetch members, ", err, c)
	}

	a.balanceRepo.Entity.MerchantID = payload.MerchantID
	a.balanceRepo.Entity.PartnerID = payload.PartnerID
	a.balanceRepo.Entity.Type = utilities.AccountTypeMerchant
	err = a.balanceRepo.GetLastBalance()
	if err != nil {
		return SendDefaultPaginationErrResponse("cannot get merchant curren balance, ", err, c)
	}

	return c.Status(200).JSON(entity.PaginatedResponseMembers{
		Success: true,
		Message: "members successfully fetched",
		Data: &entity.PaginatedResponseMemberDetails{
			Total:       total,
			LastBalance: a.balanceRepo.Entity.LastBalance,
			Result:      members,
			Pagination: entity.PaginationInfo{
				PerPage:     payload.Size,
				CurrentPage: payload.Page,
				LastPage:    pages,
			},
		},
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
			utilities.Log.Println("err: ", err2.Error())
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
		currentBalance, _ := crypt.DecryptAndConvert([]byte(account.SecretKey), account.LastBalance)

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
			utilities.Log.Println("err: ", err2.Error())
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
