package services

import (
	"context"
	"strings"
)

type EchoRequest struct {
	Input string `json:"input"`
}

type EchoResponse struct {
	Output string `json:"output"`
}

func Echo(ctx context.Context, echoRequest EchoRequest) (EchoResponse, error) {
	return EchoResponse{
		Output: strings.Replace(echoRequest.Input, "i", "o", -1),
	}, nil
}
