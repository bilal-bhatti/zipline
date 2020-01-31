package services

import (
	"strings"
)

type EchoRequest struct {
	Input string `json:"input"`
}

type EchoResponse struct {
	Output string `json:"output"`
}

func Echo(echoRequest EchoRequest) (EchoResponse, error) {
	return EchoResponse{
		Output: strings.Replace(echoRequest.Input, "i", "o", -1),
	}, nil
}
