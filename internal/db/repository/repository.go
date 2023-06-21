package repository

import (
	"github.com/dw-account-service/internal/db/entity"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	Account = AccountRepository{}
	Balance = BalanceRepository{}
	Topup   = TopupRepository{}
)

const (
	AccountTypeRegular  = 1
	AccountTypeMerchant = 2

	AccountStatusActive      = "active"
	AccountStatusDeactivated = "deactivated"

	TransSuccessStatus = "00-SUCCESS"
	TransFailedStatus  = "01-FAILED"
	TransPendingStatus = "02-PENDING"
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
