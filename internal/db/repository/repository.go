package repository

var (
	AccountRepository = Account{}
	BalanceRepository = Balance{}
	TopupRepository   = Topup{}
)

const (
	TransSuccessStatus = "00-SUCCESS"
	TransFailedStatus  = "01-FAILED"
	TransPendingStatus = "02-PENDING"
)
