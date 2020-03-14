package web

import (
	"context"
	"strings"
)

// EchoResponse model
type EchoResponse struct {
	Output string `json:"output"`
}

// Echo returns body with 'i's replaced with 'o's
func Echo(ctx context.Context, input string) (EchoResponse, error) {
	return EchoResponse{
		Output: strings.Replace(input, "i", "o", -1),
	}, nil
}
