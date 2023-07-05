package repository

import (
	"context"
	"errors"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/db/entity"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type TransactionRepository struct {
	Entity *entity.BalanceTransaction
}

func NewTransactionRepository() TransactionRepository {
	return TransactionRepository{Entity: new(entity.BalanceTransaction)}
}

// UpdateBalance :
func (t *TransactionRepository) UpdateBalance() (*entity.AccountBalance, error) {

	// 1. update balance on current document
	filter := bson.D{
		{"partnerId", t.Entity.PartnerID},
		{"merchantId", t.Entity.MerchantID},
		{"terminalId", t.Entity.TerminalID},
	}
	update := bson.D{
		{"$set", bson.D{
			{"lastBalanceNumeric", t.Entity.LastBalance},
			{"lastBalance", t.Entity.LastBalanceEncrypted},
			{"updatedAt", time.Now().UnixMilli()},
		}},
	}

	updateResult, err := db.Mongo.Collection.Account.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	if updateResult.ModifiedCount == 0 {
		return nil, errors.New("update balance failed, cannot find account with current id")
	}

	// 2. fetch current updated document
	account := new(entity.AccountBalance)
	err = db.Mongo.Collection.Account.FindOne(context.Background(), filter).Decode(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
