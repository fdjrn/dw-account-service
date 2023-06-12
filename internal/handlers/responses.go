package handlers

type Responses struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Total   int         `json:"total,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
