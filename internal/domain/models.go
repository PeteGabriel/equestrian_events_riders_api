package domain

import riders "github.com/petegabriel/hippobase"

// Competition is a struct that represents a competition inside the API.
type Competition struct {
	ID     string   `json:"competitions" jsonapi:"primary,competitions"`
	Name   string   `json:"name" jsonapi:"attr,name"`
	Events []*Event `json:"events" jsonapi:"relation,events,omitempty"`
}

// Event is a struct that represents an event inside a certain competition.
// Competitions can have zero or multiple events.
type Event struct {
	ID          string                   `json:"id" jsonapi:"primary,events"`
	Date        string                   `json:"date" jsonapi:"attr,date"`
	Name        string                   `json:"name" jsonapi:"attr,name"`
	Nations     int                      `json:"nations" jsonapi:"attr,total_of_nations"`
	Athletes    int                      `json:"athletes" jsonapi:"attr,total_of_athletes"`
	Horses      int                      `json:"horses" jsonapi:"attr,total_of_horses"`
	Competitors []*riders.RidersEntryRow `json:"competitors" jsonapi:"attr,competitors"`
}
