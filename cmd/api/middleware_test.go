package main

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func mockCompetitions() []Competition {
	var competitions []Competition
	competitions = append(competitions, Competition{
		ID:   "61. Mannheim Maimarkt Turnier",
		Name: "61. Mannheim Maimarkt Turnier",
		Events: []*Event{
			{
				ID:       "Entries CDIU25",
				Date:     "2023-10-01",
				Name:     "Event 1",
				Nations:  5,
				Athletes: 10,
				Horses:   15,
			},
			{
				ID:       "Entries CSIO3*/YH",
				Date:     "2023-10-02",
				Name:     "Event 2",
				Nations:  3,
				Athletes: 8,
				Horses:   12,
			},
		},
	})

	return competitions
}

func seedInMemoryDB() *badger.DB {
	// Initialize the in-memory database
	optionsBadger := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(optionsBadger)
	if err != nil {
		panic(err)
	}

	competitions := mockCompetitions()
	for _, c := range competitions {
		js, err := json.Marshal(&c)

		err = db.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte(c.Name), js)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	return db
}

func TestCacheCheckMiddleware(t *testing.T) {
	db := seedInMemoryDB()
	defer db.Close()

	app := &Application{
		InMemory: db,
	}

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the middleware function
	app.CheckCacheForEntryLists(c)
	// Check if the response is as expected

	// Check if the next handler was called
	if len(c.Errors) != 0 {
		t.Errorf("expected no errors, got %d", len(c.Errors))
	}

	// Check if the response status code is 200 OK
	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	var competitions []Competition
	if err := json.Unmarshal(w.Body.Bytes(), &competitions); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(competitions) < 1 {
		t.Fatalf("expected competitions, got none")
	}

}
