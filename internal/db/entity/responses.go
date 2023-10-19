package entity

type Responses struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Total   int         `json:"total,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// result total field section
// -------------------------------------------

type ResponsePayloadData struct {
	Total  int         `json:"total"`
	Result interface{} `json:"results"`
}

type ResponsePayload struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    ResponsePayloadData `json:"data,omitempty"`
}

// paginated result section
// -------------------------------------------

type PaginationInfo struct {
	PerPage     int64 `json:"perPage,omitempty"`
	CurrentPage int64 `json:"currentPage,omitempty"`
	LastPage    int64 `json:"lastPage,omitempty"`
}

type PaginatedDetailResponse struct {
	Total      int64          `json:"total,omitempty"`
	Result     interface{}    `json:"results,omitempty"`
	Pagination PaginationInfo `json:"pagination,omitempty"`
}

type PaginatedResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    PaginatedDetailResponse `json:"data,omitempty"`
}

type PaginatedResponseMemberDetails struct {
	LastBalance int64          `json:"lastBalance"`
	Total       int64          `json:"totalMember,omitempty"`
	Result      interface{}    `json:"results,omitempty"`
	Pagination  PaginationInfo `json:"pagination,omitempty"`
}

type PaginatedResponseMembers struct {
	Success bool                            `json:"success"`
	Message string                          `json:"message"`
	Data    *PaginatedResponseMemberDetails `json:"data,omitempty"`
}
