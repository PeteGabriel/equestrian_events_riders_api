package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	riders "github.com/petegabriel/equestrian_events_riders_list"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var app *Application

func init() {
	app = &Application{
		InMemory: nil, //TODO create a mock for this
	}
}

func setupMockRouter() *gin.Engine {
	//setup router to use mocked fixtures
	router := gin.Default()

	router.GET("/competitions", func(c *gin.Context) {
		f := fixtures()
		c.Writer.Header().Set("Content-Type", jsonapi.MediaType)
		c.Writer.WriteHeader(http.StatusOK)

		if err := jsonapi.MarshalPayload(c.Writer, f); err != nil {
			panic(err)
		}
	})

	return router
}

// test ListCompetitions
func TestListCompetitions(t *testing.T) {

	r := setupMockRouter()

	// create a new gin context
	w := httptest.NewRecorder()

	// call ListCompetitions
	req, _ := http.NewRequest("GET", "/competitions", nil)
	r.ServeHTTP(w, req)

	// check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var competitions []Competition

	data, err := jsonapi.UnmarshalManyPayload(w.Body, reflect.TypeOf(competitions))
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	// check the response body
	if len(data) != 2 {
		t.Errorf("expected 2 competitions, got %d", len(data))
	}
}

func fixtures() []*Competition {

	competitors := make([]riders.RidersEntryRow, 0)
	competitors = append(competitors, riders.RidersEntryRow{
		Flag:        "No country",
		CountryCode: "No code",
		CountryName: "No name",
		Pairs:       make(map[string][]string),
	})

	return []*Competition{
		{
			Name: "Test Competition #1",
			ID:   "Test Competition #1",
			Events: []*Event{
				{
					Date: "",
					Name: "Event #1",

					Nations:  1,
					Athletes: 1,
					Horses:   1,
				},
				{
					Date: "",
					Name: "Event #2",

					Nations:  1,
					Athletes: 1,
					Horses:   1,
				},
			},
		},
		{
			Name: "Test Competition #2",
			ID:   "Test Competition #2",
			Events: []*Event{
				{
					Date: "",
					Name: "Event #1",

					Nations:  1,
					Athletes: 1,
					Horses:   1,
				},
				{
					Date: "",
					Name: "Event #2",

					Nations:  1,
					Athletes: 1,
					Horses:   1,
				},
			},
		},
	}
}
