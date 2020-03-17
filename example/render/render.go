package render

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type StateError interface {
	error
	Code() int
}

type err struct {
	code int
	msg  string
}

func (e *err) Error() string {
	return e.msg
}

func NewBadRequestError(msg string) error {
	return &err{
		code: http.StatusBadRequest,
		msg:  msg,
	}
}

func NewInternalServerError(msg string) error {
	return &err{
		code: http.StatusInternalServerError,
		msg:  msg,
	}
}

func Response(w http.ResponseWriter, data interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(data); err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("failed to encode JSON response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(buf.Bytes())
	if err != nil {
		log.Println("failed to write HTTP response", err)
	}
}

type ErrorResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func toErrorResponse(err error) *ErrorResponse {
	var code int

	switch err := err.(type) {
	case StateError:
		code = err.Code()
	default:
		code = http.StatusInternalServerError
	}

	return &ErrorResponse{
		Code:   code,
		Status: http.StatusText(code),
	}
}

func Error(w http.ResponseWriter, err error) {
	log.Println(err.Error(), err)

	resp := toErrorResponse(err)

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	if err := enc.Encode(resp); err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, http.StatusText(resp.Code), resp.Code)
		log.Println("failed to encode JSON response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(resp.Code)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Println("failed to write HTTP response", err)
	}
}
