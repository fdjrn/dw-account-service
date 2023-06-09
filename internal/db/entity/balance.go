package entity

// BalanceInquiry
// adalah struct yang digunakan untuk proses pengecekan saldo akhir pengguna
type BalanceInquiry struct {
	// MongoDB ObjectID
	ID string `json:"accountId,omitempty" bson:"_id,omitempty"`

	// Unique ID (user programs id, yang didapat dari MDL MyDigiLearn yang terdaftar)
	UniqueID string `json:"uniqueId,omitempty" bson:"uniqueId,omitempty"`

	// key yg digunakan untuk proses encrypt/decrypt last balance pengguna
	SecretKey string `json:"-" bson:"secretKey,omitempty"`

	// salted/encrypted last balance pengguna
	LastBalance string `json:"-" bson:"lastBalance,omitempty"`

	// nominal last balance secara numeric
	CurrentBalance int64 `json:"currentBalance" bson:"-"`
}

// BalanceTopUp
// adalah struct yang digunakan untuk proses penambahan saldo pengguna
// dan merupakan message payload yang dikirimkan ke kafka
type BalanceTopUp struct {
	// MongoDB ObjectID
	ID string `json:"accountId,omitempty" bson:"_id,omitempty"`

	// Unique ID (user programs id, yang didapat dari MDL MyDigiLearn yang terdaftar)
	UniqueID string `json:"uniqueId" bson:"uniqueId,omitempty"`

	// Amount of topup
	Amount int `json:"topupAmount" bson:"topupAmount,omitempty"`

	// Internal Ref Number, generated by system
	InRefNumber string `json:"inRefNumber" bson:"inRefNumber,omitempty"`

	// External Ref Number, generated by third party e.g. SoF
	ExRefNumber string `json:"exRefNumber" bson:"exRefNumber,omitempty"`

	// External Transaction Date / Success Date Time (Topup),
	// generated by third party e.g. SoF
	TransDate int `json:"transDate" bson:"transDate,omitempty"`

	// Balance after addition
	LastBalance int64 `json:"currentBalance,omitempty" bson:"lastBalance,omitempty"`

	// Encrypted last balance after addition
	LastBalanceEncrypted string `json:"-" bson:"-"`

	// Receipt Number
	ReceiptNumber string `json:"receiptNumber,omitempty" bson:"receiptNumber"`

	// audit trail timestamp dalam format UNIX timestamp
	CreatedAt int64 `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt int64 `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// BalanceDeduction
// adalah struct yang digunakan untuk proses pengurangan (deduction) saldo pengguna
// dan merupakan message payload yang dikirimkan ke kafka
type BalanceDeduction struct {
	// MongoDB ObjectID
	ID string `json:"accountId,omitempty" bson:"_id,omitempty"`

	// Unique ID (user programs id, yang didapat dari MDL MyDigiLearn yang terdaftar)
	UniqueID string `json:"uniqueId"`

	// Amount of deduction
	Amount int `json:"amount"`

	// Id platform yang bekerjasama dengan wallet system (dalam hal ini MDL)
	PartnerID string `json:"partnerId,omitempty"`

	// Id merchant yang terdaftar pada platform tersebut
	// dalam hal ini adalah organisasi/instansi/perusahaan yang bekerjasama dengan platform MDL
	MerchantID int `json:"merchantID,omitempty"`

	// Tipe transaksi deduct,
	// 1: pembelian konten (default MDL) | 2: TBD
	TransType int `json:"transType"`

	// id/kode item transaksi
	ItemCode string `json:"itemCode"`

	// deskripsi transaksi / item transaksi
	Description string `json:"description"`

	// nomor invoice yg didapat dari partner
	InvoiceNumber string `json:"invoiceNumber"`

	// nomor resi/invoice yang di generate dari wallet system
	ReceiptNumber string `json:"receiptNumber,omitempty"`

	// Nilai saldo akhir pengguna setelah melakukan transaksi deduct
	LastBalance int64 `json:"currentBalance"`

	// Nilai saldo akhir yang di encrypt
	LastBalanceEncrypted string `json:"-"`
}
