package repository

import (
	"context"
	"errors"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/pkg/tools"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"time"
)

type Balance struct {
}

func (b *Balance) Inquiry(uid string) (int, entity.BalanceInquiry, error) {

	//id, _ := primitive.ObjectIDFromHex(uid)

	// filter criteria
	filter := bson.D{{"uniqueId", uid}, {"active", true}}

	var balance entity.BalanceInquiry
	err := db.Mongo.Collection.Account.FindOne(
		context.Background(),
		filter,
	).Decode(&balance)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.StatusNotFound, entity.BalanceInquiry{}, errors.New("balance not found or it has been unregistered")
		}
		return fiber.StatusInternalServerError, entity.BalanceInquiry{}, err
	}

	// convert lastBalance
	currentBalance, _ := tools.DecryptAndConvert([]byte(balance.SecretKey), balance.LastBalance)
	balance.CurrentBalance = int64(currentBalance)

	return fiber.StatusOK, balance, nil
}

// UpdateBalance is a function that update lastBalance field based on supplied uniqueId
func (b *Balance) UpdateBalance(uid string, lastBalance string) (int, error) {

	// 1. update balance on current document
	filter := bson.D{{"uniqueId", uid}}
	update := bson.D{
		{"$set", bson.D{
			{"lastBalance", lastBalance},
			{"updatedAt", time.Now().UnixMilli()},
		}},
	}

	result, err := db.Mongo.Collection.Account.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	if result.ModifiedCount == 0 {
		return fiber.StatusBadRequest,
			errors.New("update balance failed, cannot find account with current id")
	}

	return fiber.StatusOK, nil
}
