package utilities

var Charset = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const (
	TransTopUp      = "Top-Up"
	TransPayment    = "Payment"
	TransDistribute = "Distribution"

	TrxStatusSuccess          = "00"
	TrxStatusPending          = "01"
	TrxStatusPartialSuccess   = "02"
	TrxStatusInvalidParams    = "03"
	TrxStatusInvalidAccount   = "04"
	TrxStatusFailed           = "05"
	TrxStatusInsufficientFund = "06"
)
