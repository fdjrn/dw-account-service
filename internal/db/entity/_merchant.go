package entity

// MerchantTrxRequest
// adalah struct yang digunakan untuk proses penambahan saldo merchant
// dan merupakan message payload yang dikirimkan ke kafka
type MerchantTrxRequest struct {

	// Unique ID (user programs id, yang didapat dari MDL MyDigiLearn yang terdaftar)
	UniqueID      string `json:"uniqueId,omitempty"`
	PartnerID     string `json:"partnerId"`
	MerchantID    string `json:"merchantId"`
	TerminalID    string `json:"terminalId,omitempty"`
	VoucherCode   int    `json:"voucherCode,omitempty"`
	VoucherAmount int    `json:"voucherAmount,omitempty"`

	// Amount of topup
	Amount int64 `json:"topupAmount" bson:"topupAmount,omitempty"`

	// PartnerRefNumber adalah eksternal Ref Number yang didapat dari partner
	PartnerRefNumber string `json:"partnerRefNumber"`

	// PartnerTransDate adalah tgl transaksi yang dikirim oleh client/partner
	// format: YYYY-MM-DD hh:mm:ss
	PartnerTransDate string `json:"partnerTransDate"`
}
