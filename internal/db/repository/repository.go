package repository

import (
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func GetDefaultAccountFilter(account *entity.AccountBalance) bson.D {
	filter := bson.D{
		{"partnerId", account.PartnerID},
		{"merchantId", account.MerchantID},
	}

	if account.TerminalID != "" {
		filter = append(filter, bson.D{{"terminalId", account.TerminalID}}...)
	}

	if account.Type > 0 {
		filter = append(filter, bson.D{{"type", account.Type}}...)
	}

	return filter
}

func GetDefaultAccountStatusFilter(status string) bson.D {
	var filter = bson.D{}

	switch status {
	case utilities.AccountStatusActive:
		filter = bson.D{{"active", true}}
	case utilities.AccountStatusDeactivated:
		filter = bson.D{{"active", false}}
	default:
	}

	return filter
}
