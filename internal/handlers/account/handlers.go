package account

import (
	"fmt"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/internal/handlers"
	"github.com/dw-account-service/pkg/tools"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

var r = repository.AccountRepository{}

// isRegistered is a private function that check whether account id has been registered based on phoneNumber
func isRegistered(uniqueId string) bool {
	_, _, err := r.FindByUniqueID(uniqueId, true)
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
	_, _, err := r.FindByActiveStatus(uniqueId, false)
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
		return c.Status(400).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if isRegistered(a.MDLUniqueID) {
		return c.Status(400).JSON(handlers.Responses{
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

	code, id, err := r.InsertDocument(a)
	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	_, createdAccount, _ := r.FindByID(id, true)

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: "account has been successfully registered",
		Data:    createdAccount,
	})
}

// Unregister is a function that used to change active status to false (unregistered)
func Unregister(c *fiber.Ctx) error {

	// new u struct
	u := new(entity.UnregisterAccount)

	// parse body payload
	if err := c.BodyParser(u); err != nil {
		return c.Status(400).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// check if already been unregistered
	if isUnregistered(u.MDLUniqueID) {
		return c.Status(400).JSON(handlers.Responses{
			Success: false,
			Message: "account has already been unregistered",
			Data:    u,
		})
	}
	code, err := r.DeactivateAccount(u)
	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// insert into accountDeactivated collection
	code, _, err = r.InsertDeactivatedAccount(u)

	_, updatedAccount, _ := r.FindByUniqueID(u.MDLUniqueID, false)

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: "account has been successfully unregistered",
		Data:    updatedAccount,
	})
}

// Reregister is a function that used to re-activation account balance by changing active status to true
// and delete accountDeactivated collection by uniqueId.
func Reregister(c *fiber.Ctx) error {
	// TODO
	// 1. change active status to True
	// 2. remove document on accountDeactivated collection

	return nil
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
			return c.Status(fiber.StatusBadRequest).JSON(handlers.Responses{
				Success: false,
				Message: "invalid query param value, expected value is true or false",
				Data:    nil,
			})
		}
	}

	code, accounts, err := r.FindAll(queryParams)
	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	msgResponse := fmt.Sprintf("%saccounts fetched successfully ", accountStatus)

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: msgResponse,
		Data:    accounts,
	})
}

// GetRegisteredAccount is used to find registered account with active status = true
func GetRegisteredAccount(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	code, account, err := r.FindByID(id, true)

	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: "account fetched successfully ",
		Data:    account,
	})

}

func GetRegisteredAccountByUID(c *fiber.Ctx) error {
	code, account, err := r.FindByUniqueID(c.Params("uid"), true)

	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: "account fetched successfully ",
		Data:    account,
	})

}
