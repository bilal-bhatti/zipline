package models

type ThingRequest struct {
	Input string `json:"input"`
}

type ThingResponse struct {
	Output string `json:"output"`
}
