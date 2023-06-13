package request

type PaginatedAccountRequest struct {
	Status string `json:"status,omitempty"` // active | unregistered
	Page   int64  `json:"page,omitempty"`
	Size   int64  `json:"size,omitempty"`
}
