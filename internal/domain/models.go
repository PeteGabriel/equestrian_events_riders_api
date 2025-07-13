package domain

// Competition is a struct that represents a competition inside the API.
type Competition struct {
	ID           int      `json:"competitions" jsonapi:"primary,competitions"`
	Name         string   `json:"name" jsonapi:"attr,name"`
	Events       []*Event `json:"events" jsonapi:"relation,events,omitempty"`
	URL          string   `json:"homepage_url" jsonapi:"attr,homepage_url"`
	EntryListURL string   `json:"entry_list_url" jsonapi:"attr,entry_list_url,omitempty"`
}

// Event is a struct that represents an event inside a certain competition.
// Competitions can have zero or multiple events.
type Event struct {
	ID              string       `json:"id" jsonapi:"primary,events"`
	Date            string       `json:"date" jsonapi:"attr,date"`
	Name            string       `json:"name" jsonapi:"attr,name"`
	Nations         int          `json:"nations" jsonapi:"attr,total_of_nations,omitempty"`
	Athletes        int          `json:"athletes" jsonapi:"attr,total_of_athletes,omitempty"`
	Horses          int          `json:"horses" jsonapi:"attr,total_of_horses,omitempty"`
	RidersAndHorses []Competitor `json:"competitors" jsonapi:"attr,competitors,omitempty"`
}

type Competitor struct {
	Rider  string   `json:"rider" jsonapi:"attr,rider"`
	Horses []string `json:"horses" jsonapi:"attr,horses"`
}
