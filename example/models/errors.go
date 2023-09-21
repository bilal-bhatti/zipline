package models

type ErrorResponse struct {
	Code        int    `json:"code"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
}
