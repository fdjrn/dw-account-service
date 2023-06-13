package handlers

import (
	"fmt"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/pkg/payload/request"
	"github.com/dw-account-service/pkg/tools"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

// isRegistered is a private function that check whether account id has been registered based on phoneNumber
func isRegistered(uniqueId string) bool {
	_, _, err := repository.AccountRepository.FindByUniqueID(uniqueId, true)
	if err != nil {
		// no document found, its mean it can be registered
		if err == mongo.ErrNoDocuments {
			return false
		}

		// TODO handling unknown error
		return true
	}
	return true
}

// isUnregistered is a private function that check whether account id has been unregistered or not
func isUnregistered(uniqueId string) bool {
	_, _, err := repository.AccountRepository.FindByActiveStatus(uniqueId, false)
	if err != nil {
		// no document found, its mean it can be unregistered
		if err == mongo.ErrNoDocuments {
			return false
		}

		// TODO handling unknown error
		return true
	}
	return true
}

// Register is a function that used to insert new document into collection and set active status to true.
func Register(c *fiber.Ctx) error {

	// new account struct
	a := new(entity.AccountBalance)

	// parse body payload
	if err := c.BodyParser(a); err != nil {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if isRegistered(a.UniqueID) {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: "uniqueId has already been registered",
			Data:    a,
		})
	}

	// set default value for accountBalance document
	key, _ := tools.GenerateSecretKey()
	encryptedBalance, _ := tools.Encrypt([]byte(key), fmt.Sprintf("%016s", "0"))

	a.ID = ""
	a.SecretKey = key
	a.Active = true
	a.LastBalance = encryptedBalance
	a.MainAccountID = "-"
	a.CreatedAt = time.Now().UnixMilli()
	a.UpdatedAt = a.CreatedAt

	code, id, err := repository.AccountRepository.InsertDocument(a)
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
		Message: "account successfully registered",
		Data:    createdAccount,
	})
}

// Unregister is a function that used to change active status to false (unregistered)
func Unregister(c *fiber.Ctx) error {

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
	if isUnregistered(u.UniqueID) {
		return c.Status(400).JSON(Responses{
			Success: false,
			Message: "account has already been unregistered",
			Data:    u,
		})
	}
	code, err := repository.AccountRepository.DeactivateAccount(u)
	if err != nil {
		return c.Status(code).JSON(Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// insert into accountDeactivated collection
	code, _, err = repository.AccountRepository.InsertDeactivatedAccount(u)

	_, updatedAccount, _ := repository.AccountRepository.FindByUniqueID(u.UniqueID, false)

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: "account successfully unregistered",
		Data:    updatedAccount,
	})
}

// GetAllRegisteredAccount is used to find all registered account and can be filtered with their active status
func GetAllRegisteredAccount(c *fiber.Ctx) error {
	accountStatus := ""
	queryParams := c.Query("active")
	if queryParams != "" {
		switch strings.ToLower(queryParams) {
		case "true":
			accountStatus = "active "
		case "false":
			accountStatus = "unregistered "
		default:
			return c.Status(fiber.StatusBadRequest).JSON(Responses{
				Success: false,
				Message: "invalid query param value, expected value is true or false",
				Data:    nil,
			})
		}
	}

	code, accounts, count, err := repository.AccountRepository.FindAll(queryParams)
	if err != nil {
		return c.Status(code).JSON(ResponsePayload{
			Success: false,
			Message: err.Error(),
			Data: ResponsePayloadData{
				Total:  0,
				Result: nil,
			},
		})
	}

	msgResponse := fmt.Sprintf("%saccounts fetched successfully ", accountStatus)

	return c.Status(code).JSON(Responses{
		Success: true,
		Message: msgResponse,
		Data: ResponsePayloadData{
			Total:  count,
			Result: accounts,
		},
	})
}

func GetAllRegisteredAccountPaginated(c *fiber.Ctx) error {

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

	validStatus := map[string]interface{}{"all": 0, "active": 1, "unregistered": 2}
	if req.Status != "" {
		if _, ok := validStatus[req.Status]; !ok {
			return c.Status(400).JSON(Responses{
				Success: false,
				Message: "invalid status value. its only accept all, active or unregistered",
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

	code, accounts, total, pages, err := repository.AccountRepository.FindAllPaginated(req)
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

// GetRegisteredAccount is used to find registered account with active status = true
func GetRegisteredAccount(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	code, account, err := repository.AccountRepository.FindByID(id, true)

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

// GetRegisteredAccountByUID is used to find registered account based on uniqueId
func GetRegisteredAccountByUID(c *fiber.Ctx) error {
	code, account, err := repository.AccountRepository.FindByActiveStatus(c.Params("uid"), true)

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
