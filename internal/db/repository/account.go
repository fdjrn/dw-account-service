package repository

import (
	"context"
	"errors"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/pkg/payload/request"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

type Account struct {
}

func (a *Account) findByFilter(filter interface{}) (interface{}, error) {
	account := new(entity.AccountBalance)

	// filter condition
	err := db.Mongo.Collection.Account.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

/*
FindAll
function args:

	queryParams: string

return:

	code: int,
	accounts: interface{},
	length: int,
	err: error
*/
func (a *Account) FindAll(queryParams string) (int, interface{}, int, error) {
	filter := bson.D{}

	if queryParams != "" {
		filter = bson.D{{"active", false}}
		if queryParams == "true" {
			filter = bson.D{{"active", true}}
		}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()

	cursor, err := db.Mongo.Collection.Account.Find(ctx, filter,
		options.Find().SetProjection(bson.D{{"secretKey", 0}, {"lastBalance", 0}}),
	)
	if err != nil {
		return fiber.StatusInternalServerError, nil, 0, err
	}

	var accounts []entity.AccountBalance
	if err = cursor.All(context.TODO(), &accounts); err != nil {
		return fiber.StatusInternalServerError, nil, 0, err
	}

	return fiber.StatusOK, &accounts, len(accounts), nil
}

/*
FindAllPaginated
function args:

	*request.PaginatedAccountRequest

return:

	code int,
	data interfaces{},
	totalDocument int64,
	totalPages int,
	err error
*/
func (a *Account) FindAllPaginated(request *request.PaginatedAccountRequest) (int, interface{}, int64, int64, error) {
	var filter interface{}
	switch request.Status {
	case "active":
		filter = bson.D{{"active", true}}
	case "unregistered":
		filter = bson.D{{"active", false}}
	default:
		filter = bson.D{}
	}

	skipValue := (request.Page - 1) * request.Size

	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()

	cursor, err := db.Mongo.Collection.Account.Find(
		ctx,
		filter,
		options.Find().
			SetProjection(bson.D{{"secretKey", 0}, {"lastBalance", 0}}).
			SetSkip(skipValue).
			SetLimit(request.Size),
	)

	if err != nil {
		return fiber.StatusInternalServerError, nil, 0, 0, err
	}

	totalDocs, _ := db.Mongo.Collection.Account.CountDocuments(ctx, filter)
	var accounts []entity.AccountBalance
	if err = cursor.All(context.TODO(), &accounts); err != nil {
		return fiber.StatusInternalServerError, nil, 0, 0, err
	}

	if len(accounts) == 0 {
		return fiber.StatusInternalServerError, nil, 0, 0, errors.New("empty results or last pages has been reached")
	}

	totalPages := math.Ceil(float64(totalDocs) / float64(request.Size))
	return fiber.StatusOK, &accounts, totalDocs, int64(totalPages), nil
}

// FindByID : id args accept interface{} or primitive.ObjectID
// make sure to convert it first
func (a *Account) FindByID(id interface{}, active bool) (int, interface{}, error) {
	// filter condition
	filter := bson.D{{"_id", id}}
	if active {
		filter = bson.D{{"_id", id}, {"active", true}}
	}

	account, err := a.findByFilter(filter)
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, account, nil
}

func (a *Account) FindByUniqueID(id string, active bool) (int, interface{}, error) {

	filter := bson.D{{"uniqueId", id}}
	if active {
		filter = bson.D{{"uniqueId", id}, {"active", true}}
	}

	account, err := a.findByFilter(filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			//return fiber.StatusNotFound, nil, errors.New("account not found or it has been unregistered")
			return fiber.StatusNotFound, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, account, nil
}

func (a *Account) FindByActiveStatus(id string, status bool) (int, interface{}, error) {
	//id, _ := primitive.ObjectIDFromHex(accountId)

	account, err := a.findByFilter(bson.D{
		{Key: "uniqueId", Value: id},
		{Key: "active", Value: status},
	})

	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, account, nil
}

func (a *Account) DeactivateAccount(u *entity.UnregisterAccount) (int, error) {
	//id, _ := primitive.ObjectIDFromHex(u.UniqueID)

	// filter condition
	filter := bson.D{{"uniqueId", u.UniqueID}}

	// update field
	update := bson.D{
		{"$set", bson.D{
			{"active", false},
			{"updatedAt", time.Now().UnixMilli()},
		}},
	}

	result, err := db.Mongo.Collection.Account.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	if result.ModifiedCount == 0 {
		return fiber.StatusBadRequest, errors.New(
			"update failed, cannot find account with current uniqueId")
	}

	return fiber.StatusOK, nil
}

func (a *Account) InsertDocument(account *entity.AccountBalance) (int, interface{}, error) {

	result, err := db.Mongo.Collection.Account.InsertOne(context.TODO(), account)

	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusCreated, result.InsertedID, nil

}

func (a *Account) InsertDeactivatedAccount(account *entity.UnregisterAccount) (int, interface{}, error) {

	result, err := db.Mongo.Collection.UnregisterAccount.InsertOne(context.TODO(), account)

	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusCreated, result.InsertedID, nil

}

func (a *Account) RemoveDeactivatedAccount(acc *entity.UnregisterAccount) (int, error) {
	filter := bson.D{{"uniqueId", acc.UniqueID}}
	result, err := db.Mongo.Collection.UnregisterAccount.DeleteOne(context.TODO(), filter)

	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	if result.DeletedCount > 0 {
		return fiber.StatusNoContent, nil
	}

	return fiber.StatusInternalServerError, errors.New("remove deactivated account failed, no document found")
}
