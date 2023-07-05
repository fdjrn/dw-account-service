package entity

import "time"

type PeriodsRequest struct {
	Start     string    `json:"start,omitempty"`
	StartDate time.Time `json:"-"`
	End       string    `json:"end,omitempty"`
	EndDate   time.Time `json:"-"`
}

type PaginatedAccountRequest struct {
	PartnerID  string         `json:"partnerId,omitempty"`
	MerchantID string         `json:"merchantID,omitempty"`
	Type       int            `json:"type,omitempty"`
	Status     string         `json:"status,omitempty"` // { all | active | deactivated }
	Periods    PeriodsRequest `json:"periods,omitempty"`
	Page       int64          `json:"page,omitempty"`
	Size       int64          `json:"size,omitempty"`
}
