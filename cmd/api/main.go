package main

import "github.com/dgraph-io/badger/v4"

func main() {

	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)

	app := Application{
		InMemory: db,
	}

	err = app.serve()
	if err != nil {
		panic(err)
	}
}
