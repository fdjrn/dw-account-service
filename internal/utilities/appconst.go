package utilities

var Charset = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const (
	AccountTypeRegular  = 1
	AccountTypeMerchant = 2

	AccountStatusActive      = "active"
	AccountStatusDeactivated = "deactivated"
	AccountStatusAll         = "all"

	TransTypeTopUp        = 1 //"Top-Up"
	TransTypePayment      = 2 //"Payment"
	TransTypeDistribution = 3 //"Distribution"

	TrxStatusSuccess          = "00"
	TrxStatusPending          = "01"
	TrxStatusPartialSuccess   = "02"
	TrxStatusInvalidParams    = "03"
	TrxStatusInvalidAccount   = "04"
	TrxStatusFailed           = "05"
	TrxStatusInsufficientFund = "06"
	TrxStatusCallbackFailed   = "07"
)
