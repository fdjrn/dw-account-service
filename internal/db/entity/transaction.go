package entity

type TransactionItem struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Code   string `json:"code,omitempty" bson:"code"`
	Name   string `json:"name" bson:"name"`
	Amount int64  `json:"amount" bson:"amount"`
	Price  int64  `json:"price,omitempty" bson:"price"`
	Qty    int    `json:"qty,omitempty" bson:"qty"`
}

type BalanceTransaction struct {
	ID                   string            `json:"id,omitempty" bson:"_id,omitempty"`
	TransDate            string            `json:"transDate,omitempty" bson:"transDate"`               // YYYY-MM-DD hh:mm:ss
	TransDateNumeric     int64             `json:"transDateNumeric,omitempty" bson:"transDateNumeric"` // unix time millis
	ReferenceNo          string            `json:"referenceNo,omitempty" bson:"referenceNo"`
	ReceiptNumber        string            `json:"receiptNumber,omitempty" bson:"receiptNumber"`
	LastBalance          int64             `json:"lastBalance,omitempty" bson:"lastBalance"`
	LastBalanceEncrypted string            `json:"-" bson:"-"`
	Status               string            `json:"status,omitempty" bson:"status"`
	TransType            int               `json:"transType,omitempty" bson:"transType"` // (1) TopUp | (2) Payment | (3) Distribution
	PartnerTransDate     string            `json:"partnerTransDate" bson:"partnerTransDate"`
	PartnerRefNumber     string            `json:"partnerRefNumber" bson:"partnerRefNumber"`
	PartnerID            string            `json:"partnerId" bson:"partnerId"`
	MerchantID           string            `json:"merchantId" bson:"merchantId"`
	TerminalID           string            `json:"terminalId" bson:"terminalId"`
	TerminalName         string            `json:"terminalName" bson:"terminalName"`
	TotalAmount          int64             `json:"totalAmount" bson:"totalAmount"`
	Items                []TransactionItem `json:"items" bson:"items"`
	CreatedAt            int64             `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt            int64             `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	RequestDetail        RequestDetail     `json:"requestDetail" bson:"requestDetail"`
}

type BalanceDistributionInfo struct {
	Index       int
	Data        BalanceTransaction
	WorkerIndex int
	Err         error
}

type RequestDetail struct {
	Origin    string `json:"origin" bson:"origin"`
	Timestamp string `json:"timestamp,omitempty" bson:"timestamp"`
}
