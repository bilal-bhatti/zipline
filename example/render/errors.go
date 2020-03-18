package render

import (
	"net/http"

	"github.com/pkg/errors"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type StateError interface {
	error
	Code() int
	StackTracer
}

type stateError struct {
	code int
	err  error
}

func (e *stateError) Code() int {
	return e.code
}

func (e *stateError) Error() string {
	return e.err.Error()
}

func (e *stateError) StackTrace() errors.StackTrace {
	if err, ok := e.err.(StackTracer); ok {
		return err.StackTrace()
	}

	return errors.StackTrace{}
}

func WrapBadRequestError(err error, msg string) error {
	return wrapError(http.StatusBadRequest, err, msg)
}

func NewBadRequestError(msg string) error {
	return newError(http.StatusBadRequest, msg)
}

func WrapNotFoundError(err error, msg string) error {
	return wrapError(http.StatusNotFound, err, msg)
}

func NewNotFoundError(msg string) error {
	return newError(http.StatusNotFound, msg)
}

func WrapForbiddenError(err error, msg string) error {
	return wrapError(http.StatusForbidden, err, msg)
}

func NewForbiddenError(msg string) error {
	return newError(http.StatusForbidden, msg)
}

func wrapError(code int, err error, msg string) StateError {
	return &stateError{
		code: code,
		err:  errors.Wrap(err, msg),
	}
}

func newError(code int, msg string) StateError {
	return &stateError{
		code: code,
		err:  errors.New(msg),
	}
}
