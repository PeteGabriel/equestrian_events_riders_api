package main

import (
	"github.com/dgraph-io/badger/v4"
)

type Application struct {
	InMemory *badger.DB
}
