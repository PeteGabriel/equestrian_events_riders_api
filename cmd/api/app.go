package main

import (
	"errors"
	"github.com/dgraph-io/badger/v4"
	"net/http"
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

func InternalError(message string) *HTTPError {
	return New(http.StatusInternalServerError, message)
}
