package main

import (
	"equestrian-events-api/internal/application"

	"github.com/dgraph-io/badger/v4"
)

func main() {

	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)

	app := application.Application{
		InMemory: db,
	}

	err = app.Serve()
	if err != nil {
		panic(err)
	}
}
