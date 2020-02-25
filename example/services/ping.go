package services

import (
	"context"
	"strings"
)

type PingRequest struct {
	Input string `json:"input"`
}

type PingResponse struct {
	Output string `json:"output"`
}

// Ping returns body with 'i's replaced with 'o's
func Ping(ctx context.Context, pingRequest PingRequest) (PingResponse, error) {
	return PingResponse{
		Output: strings.Replace(pingRequest.Input, "i", "o", -1),
	}, nil
}
