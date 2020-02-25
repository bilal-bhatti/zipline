package web

import (
	"context"
	"strings"
)

// EchoRequest model
type EchoRequest struct {
	Input string `json:"input"`
}

// EchoResponse model
type EchoResponse struct {
	Output string `json:"output"`
}

// Echo returns body with 'i's replaced with 'o's
func Echo(ctx context.Context, echoRequest EchoRequest) (EchoResponse, error) {
	return EchoResponse{
		Output: strings.Replace(echoRequest.Input, "i", "o", -1),
	}, nil
}
