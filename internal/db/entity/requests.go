package entity

type PaginatedAccountRequest struct {
	PartnerID string `json:"partnerId,omitempty"`

	MerchantID string `json:"merchantID,omitempty"`

	Type int `json:"type,omitempty"`

	// Account status: active/deactivated
	Status string `json:"status,omitempty"`

	Page int64 `json:"page,omitempty"`

	Size int64 `json:"size,omitempty"`
}
