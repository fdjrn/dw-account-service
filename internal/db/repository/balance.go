package repository

import (
	"context"
	"errors"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/utilities/crypt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type BalanceRepository struct {
	Entity *entity.InquiryBalance
	//Inquiry *entity.InquiryBalance
}

func NewBalanceRepository() BalanceRepository {
	return BalanceRepository{
		Entity: new(entity.InquiryBalance),
		//Inquiry: new(entity.InquiryBalance),
	}
}

func (b *BalanceRepository) GetLastBalance() error {
	// filter criteria
	filter := bson.D{
		{"active", true},
		{"partnerId", b.Entity.PartnerID},
		{"merchantId", b.Entity.MerchantID},
		{"type", b.Entity.Type},
	}

	if b.Entity.TerminalID != "" {
		filter = append(filter, bson.D{{"terminalId", b.Entity.TerminalID}}...)
	}

	err := db.Mongo.Collection.Account.FindOne(
		context.Background(),
		filter,
		options.FindOne().SetProjection(bson.D{
			{"_id", 0},
			//{"uniqueId", 0},
			//{"active", 0},
			//{"type", 0},
			//{"createdAt", 0},
			//{"updatedAt", 0},
			{"partnerId", 1},
			{"merchantId", 1},
			{"terminalId", 1},
			{"lastBalanceNumeric", 1},
		}),
	).Decode(b.Entity)

	if err != nil {
		return err
	}

	return nil
}

func (b *BalanceRepository) InquiryBalance(uid string) (int, entity.BalanceInquiry, error) {

	//id, _ := primitive.ObjectIDFromHex(uid)

	// filter criteria
	filter := bson.D{{"uniqueId", uid}, {"active", true}}

	var balance entity.BalanceInquiry
	err := db.Mongo.Collection.Account.FindOne(
		context.Background(),
		filter,
	).Decode(&balance)

	if err != nil {
		return fiber.StatusInternalServerError, entity.BalanceInquiry{}, err
	}

	// convert lastBalance
	currentBalance, _ := crypt.DecryptAndConvert([]byte(balance.SecretKey), balance.LastBalance)
	balance.CurrentBalance = int64(currentBalance)

	return fiber.StatusOK, balance, nil
}

func (b *BalanceRepository) MerchantInquiryBalance(inquiry entity.BalanceInquiry) (int, entity.BalanceInquiry, error) {

	// filter criteria
	filter := bson.D{
		{"partnerId", inquiry.PartnerID},
		{"merchantId", inquiry.MerchantID},
	}

	var balance entity.BalanceInquiry
	err := db.Mongo.Collection.Account.FindOne(
		context.Background(),
		filter,
		options.FindOne().SetProjection(bson.D{{"_id", 0}}),
	).Decode(&balance)

	if err != nil {
		return fiber.StatusInternalServerError, entity.BalanceInquiry{}, err
	}

	// convert lastBalance
	//currentBalance, _ := tools.DecryptAndConvert([]byte(balance.SecretKey), balance.LastBalance)
	//balance.CurrentBalance = int64(currentBalance)

	return fiber.StatusOK, balance, nil
}

// UpdateBalance is a function that update lastBalance field based on supplied uniqueId
func (b *BalanceRepository) UpdateBalance(uid string, lastBalance string) (int, error) {

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

func (b *BalanceRepository) UpdateMerchantBalance(t *entity.BalanceTopUp) (int, error) {

	// 1. update balance on current document
	filter := bson.D{{"partnerId", t.PartnerID}, {"merchantId", t.MerchantID}}
	update := bson.D{
		{"$set", bson.D{
			{"lastBalanceNumeric", t.LastBalance},
			{"lastBalance", t.LastBalanceEncrypted},
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
