package entity

type AccountBalance struct {
	/*
		ID				| Akun (Wallet) ID yang di-generate ketika berhasil melakukan registrasi
		UniqueID 	| Unique ID yang didapat dari MDL MyDigiLearn yang terdaftar
		SecretKey 		| Key untuk melakukan proses encrypt dan decrypt lastBalance, yang di-generate ketika registrasi
		Active 			| Status Account Balance (wallet) pengguna. Value -->> active: true/false
		Type 			| Tipe Wallet pengguna, expected value -->> 1:Regular Account, 2: Verified Account
		LastBalance 	| Hashed/Encrypted nilai saldo akhir (lastBalance)
		MainAccountID	| Akun (Wallet ID) utama jika kedepannya setiap akun bisa memiliki akun turunan
	*/

	ID            string `json:"accountId,omitempty" bson:"_id,omitempty"`
	UniqueID      string `json:"uniqueId,omitempty" bson:"uniqueId,omitempty"`
	SecretKey     string `json:"-" bson:"secretKey,omitempty"`
	Active        bool   `json:"active" bson:"active"`
	Type          int    `json:"type,omitempty" bson:"type,omitempty"`
	LastBalance   string `json:"-" bson:"lastBalance"`
	MainAccountID string `json:"mainAccountID,omitempty" bson:"mainAccountID,omitempty"`
	CreatedAt     int64  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt     int64  `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type UnregisterAccount struct {
	//ID                string `json:"accountId,omitempty" bson:"_id,omitempty"`
	UniqueID          string `json:"uniqueId" bson:"uniqueId"`
	ReasonCode        int    `json:"reasonCode" bson:"reasonCode"`
	ReasonDescription string `json:"reasonDescription" bson:"reasonDescription"`
}
