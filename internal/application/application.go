package application

import (
	"errors"
	"net/http"

	"github.com/dgraph-io/badger/v4"
)

type Application struct {
	InMemory *badger.DB
}

type HTTPError struct {
	error
	Code int
}

func New(code int, message string) *HTTPError {
	return &HTTPError{
		error: errors.New(message),
		Code:  code,
	}
}

func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, message)
}

func InternalError(message string) *HTTPError {
	return New(http.StatusInternalServerError, message)
}
