package services

import (
	"context"
	"net/http"
)

func ProvideContext(r *http.Request) context.Context {
	return r.Context()
}
