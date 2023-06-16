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

	// Tipe Wallet pengguna, expected value -->> 1:Regular Account, 2: Subordinate Account
	Type int `json:"type" bson:"type"`

	// Hashed/Encrypted nilai saldo akhir (lastBalance)
	LastBalance string `json:"-" bson:"lastBalance"`

	// Akun utama jika kedepannya setiap akun bisa memiliki akun turunan
	// value dari field ini adalah UniqueID yang menjadi
	MainAccountID string `json:"mainAccountID" bson:"mainAccountID"`

	// Id platform yang bekerjasama dengan wallet system (dalam hal ini MDL)
	PartnerID string `json:"partnerId" bson:"partnerId,omitempty"`

	// Id merchant yang terdaftar pada platform tersebut
	// dalam hal ini adalah organisasi/instansi/perusahaan yang bekerjasama dengan platform MDL
	MerchantID string `json:"merchantId" bson:"merchantId,omitempty"`

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
