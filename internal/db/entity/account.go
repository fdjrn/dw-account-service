package entity

type AccountBalance struct {

	// Akun (Wallet) ID yang di-generate ketika berhasil melakukan registrasi
	ID string `json:"accountId,omitempty" bson:"_id,omitempty"`

	// Unique ID yang didapat dari MDL MyDigiLearn yang terdaftar
	UniqueID string `json:"uniqueId" bson:"uniqueId"`

	// Key untuk melakukan proses encrypt dan decrypt lastBalance, yang di-generate ketika registrasi
	SecretKey string `json:"-" bson:"secretKey"`

	// Status Account Balance (wallet) pengguna. Value -->> active: true/false
	Active bool `json:"active" bson:"active"`

	// Tipe Wallet pengguna, expected value -->> 1:Regular Account, 2: Admin Account
	Type int `json:"type" bson:"type"`

	// Hashed/Encrypted nilai saldo akhir (lastBalance)
	LastBalance string `json:"-" bson:"lastBalance"`

	// saldo akhir secara numeric
	LastBalanceNumeric int64 `json:"lastBalance" bson:"lastBalanceNumeric"`

	// Akun utama jika kedepannya setiap akun bisa memiliki akun turunan
	// value dari field ini adalah UniqueID yang menjadi
	MainAccountID string `json:"mainAccountID,omitempty" bson:"mainAccountID,omitempty"`

	// Id platform yang bekerjasama dengan wallet system (dalam hal ini MDL)
	PartnerID string `json:"partnerId,omitempty" bson:"partnerId"`

	// Id merchant yang terdaftar pada platform tersebut
	// dalam hal ini adalah organisasi/instansi/perusahaan yang bekerjasama dengan platform MDL
	MerchantID string `json:"merchantId,omitempty" bson:"merchantId"`

	// TerminalID adalah unique id yang didapat dari client,
	// bisa berupa user id, dan atau yang lainnya
	// yang nantinya dijadikan sebagai pencarian dan proses transaksi lainnya
	TerminalID string `json:"terminalId,omitempty" bson:"terminalId"`

	// TerminalName adalah deskripsi dari terminal id yang di kirim,
	// field ini bersifat optional
	TerminalName string `json:"terminalName,omitempty" bson:"terminalName,omitempty"`

	// audit trail dalam format UNIX timestamp
	CreatedAt int64 `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt int64 `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type UnregisterAccount struct {
	//ID                string `json:"accountId,omitempty" bson:"_id,omitempty"`
	UniqueID          string `json:"uniqueId" bson:"uniqueId"`
	ReasonCode        int    `json:"reasonCode" bson:"reasonCode"`
	ReasonDescription string `json:"reasonDescription" bson:"reasonDescription"`
}
