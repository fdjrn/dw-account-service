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

type ResponsePayloadDataPaginated struct {
	Result      interface{} `json:"results,omitempty"`
	Total       int64       `json:"total,omitempty"`
	PerPage     int64       `json:"perPage,omitempty"`
	CurrentPage int64       `json:"currentPage,omitempty"`
	LastPage    int64       `json:"lastPage,omitempty"`
}

type ResponsePayloadPaginated struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    ResponsePayloadDataPaginated `json:"data,omitempty"`
}
