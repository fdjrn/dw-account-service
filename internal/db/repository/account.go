package repository

import (
	"context"
	"errors"
	"github.com/dw-account-service/internal/conn/mongodb"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type AccountRepository struct {
}

func (a *AccountRepository) findByFilter(filter interface{}) (interface{}, error) {
	account := new(entity.AccountBalance)

	// filter condition
	err := mongodb.Collection.Account.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return account, nil
}

func (a *AccountRepository) FindAll(queryParams string) (int, interface{}, error) {
	filter := bson.D{}

	if queryParams != "" {
		filter = bson.D{{"active", false}}
		if queryParams == "true" {
			filter = bson.D{{"active", true}}
		}
	}

	cursor, err := mongodb.Collection.Account.Find(
		context.TODO(), filter,
		options.Find().SetProjection(bson.D{{"secretKey", 0}, {"lastBalance", 0}}),
	)
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	var accounts []entity.AccountBalance
	if err = cursor.All(context.TODO(), &accounts); err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, &accounts, nil
}

// FindByID : id args accept interface{} or primitive.ObjectID
// make sure to convert it first
func (a *AccountRepository) FindByID(id interface{}, active bool) (int, interface{}, error) {
	// filter condition
	filter := bson.D{{"_id", id}}
	if active {
		filter = bson.D{{"_id", id}, {"active", true}}
	}

	account, err := a.findByFilter(filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.StatusNotFound, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, account, nil
}

func (a *AccountRepository) FindByUniqueID(id string, active bool) (int, interface{}, error) {

	filter := bson.D{{"uniqueId", id}}
	if active {
		filter = bson.D{{"uniqueId", id}, {"active", true}}
	}

	account, err := a.findByFilter(filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.StatusNotFound, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, account, nil
}

func (a *AccountRepository) FindByActiveStatus(id string, status bool) (int, interface{}, error) {
	//id, _ := primitive.ObjectIDFromHex(accountId)

	account, err := a.findByFilter(bson.D{
		//{Key: "_id", Value: id},
		{Key: "uniqueId", Value: id},
		{Key: "active", Value: status},
	})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.StatusOK, nil, err
		}
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, account, nil
}

func (a *AccountRepository) DeactivateAccount(u *entity.UnregisterAccount) (int, error) {
	//id, _ := primitive.ObjectIDFromHex(u.MDLUniqueID)

	// filter condition
	filter := bson.D{{"uniqueId", u.MDLUniqueID}}

	// update field
	update := bson.D{
		{"$set", bson.D{
			{"active", false},
			{"updatedAt", time.Now().UnixMilli()},
		}},
	}

	result, err := mongodb.Collection.Account.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	if result.ModifiedCount == 0 {
		return fiber.StatusBadRequest, errors.New(
			"update failed, cannot find account with current uniqueId")
	}

	return fiber.StatusOK, nil
}

func (a *AccountRepository) InsertDocument(account *entity.AccountBalance) (int, interface{}, error) {

	result, err := mongodb.Collection.Account.InsertOne(context.TODO(), account)

	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusCreated, result.InsertedID, nil

}

func (a *AccountRepository) InsertDeactivatedAccount(account *entity.UnregisterAccount) (int, interface{}, error) {

	result, err := mongodb.Collection.UnregisterAccount.InsertOne(context.TODO(), account)

	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusCreated, result.InsertedID, nil

}

func (a *AccountRepository) RemoveDeactivatedAccount(acc *entity.UnregisterAccount) (int, error) {
	filter := bson.D{{"uniqueId", acc.MDLUniqueID}}
	result, err := mongodb.Collection.UnregisterAccount.DeleteOne(context.TODO(), filter)

	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	if result.DeletedCount > 0 {
		return fiber.StatusNoContent, nil
	}

	return fiber.StatusInternalServerError, errors.New("remove deactivated account failed, no document found")
}
