package repository

import (
	"context"
	"github.com/dw-account-service/internal/conn/mongodb"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Topup struct {
}

func (t *Topup) CreateTopupDocument(e *entity.BalanceTopUp) (int, error) {
	_, err := mongodb.Collection.BalanceTopup.InsertOne(context.TODO(), e)

	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusCreated, nil
}

func (t *Topup) IsUsedExRefNumber(refNo string) bool {

	// filter condition
	filter := bson.D{{"exRefNumber", refNo}}

	topup := new(entity.BalanceTopUp)
	err := mongodb.Collection.BalanceTopup.FindOne(context.TODO(), filter).Decode(&topup)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		return true
	}

	return true

}
