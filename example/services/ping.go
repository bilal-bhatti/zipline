package services

import (
	"context"
	"strings"

	"github.com/bilal-bhatti/zipline/example/connectors"
)

type PingRequest struct {
	Input string `json:"input"`
}

type PingResponse struct {
	Output string `json:"output"`
}

// Ping returns body with 'i's replaced with 'o's
func Ping(ctx context.Context, env *connectors.Env, pingRequest PingRequest) (PingResponse, error) {
	return PingResponse{
		Output: strings.Replace(pingRequest.Input, "i", "o", -1),
	}, nil
}
