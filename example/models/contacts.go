package models

type ContactRequest struct {
	Input string `json:"input"`
}

type ContactResponse struct {
	Output string `json:"output"`
}
