package main

// Competition is a struct that represents a competition inside the API.
type Competition struct {
	ID     string   `jsonapi:"primary,competitions"`
	Name   string   `jsonapi:"attr,name"`
	Events []*Event `jsonapi:"relation,events,omitempty"`
}

// Event is a struct that represents an event inside a certain competition.
// Competitions can have zero or multiple events.
type Event struct {
	ID       string `jsonapi:"primary,events"`
	Date     string `jsonapi:"attr,date"`
	Name     string `jsonapi:"attr,name"`
	Nations  int    `jsonapi:"attr,total_of_nations"`
	Athletes int    `jsonapi:"attr,total_of_athletes"`
	Horses   int    `jsonapi:"attr,total_of_horses"`
	//Competitors []riders.RidersEntryRow `jsonapi:"attr,competitors"`
}
