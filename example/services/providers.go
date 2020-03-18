package services

import (
	"context"
	"net/http"
	"net/url"
)

func ProvideContext(r *http.Request) context.Context {
	return r.Context()
}

func ProvideForwadedHeader(req *http.Request) *url.URL {
	var u *url.URL
	if xForwarded := req.Header.Get("X-Forwarded"); xForwarded != "" {
		u, _ = url.Parse("https://" + xForwarded)
	}
	return u
}
