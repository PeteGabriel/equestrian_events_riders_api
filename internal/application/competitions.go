package application

import (
	"encoding/json"
	"equestrian-events-api/internal/domain"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/petegabriel/hippobase"
	"net/http"
	"time"
)

// ListCompetitions is a handler that returns a list of competitions
// from the in-memory cache or from the module
// if the cache is empty
// @Summary List competitions
// @Description List competitions
// @Produce json
// @Success 200 {array} CompetitionList
// @Router /competitions [get]
func (app *Application) ListCompetitions() (CompetitionList, error) {

	events, err := hippobase.GetEvents()
	if err != nil {
		// TODO log the error
		return nil, New(http.StatusServiceUnavailable, "Unable to fetch events from Hippobase")
	}

	var competitions CompetitionList

	for _, parsed := range events {

		comp := &domain.Competition{
			ID:     uuid.New().String(),
			Name:   fmt.Sprintf("%s - %s", parsed.Location, parsed.Name),
			URL:    parsed.EventURL,
			Events: make([]*domain.Event, 0),
		}

		//TODO get entrylist for competition
		if parsed.EntryListURL != "" {

			eventWithEntryList, err := hippobase.GetEntryLists(parsed.EntryListURL)
			if err != nil {
				return nil, New(http.StatusServiceUnavailable,
					fmt.Sprintf("Unable to fetch entry lists from Hippobase - %s", err.Error()))
			}

			for _, evt := range eventWithEntryList.Events {
				event := &domain.Event{
					ID:       uuid.New().String(),
					Date:     evt.CreatedAt,
					Name:     evt.EventFullName,
					Nations:  evt.TotalNations,
					Athletes: evt.TotalAthletes,
					Horses:   evt.TotalHorses,
				}

				comp.Events = append(comp.Events, event)
			}
		}

		competitions = append(competitions, comp)
	}

	//TODO explore the possibility of doing this in a different goroutine
	if err := cacheEvents(competitions, app); err != nil {
		// TODO log the error
		return nil, InternalError(err.Error())
	}

	return competitions, nil
}

func cacheEvents(list CompetitionList, a *Application) error {
	if a.InMemory == nil {
		return nil
	}
	for _, cpt := range list {
		err := a.InMemory.Update(func(txn *badger.Txn) error {
			mEvt, err := json.Marshal(cpt)
			if err != nil {
				return err
			}

			newEntry := badger.
				NewEntry([]byte(cpt.Name), mEvt).
				WithTTL(time.Hour)
			err = txn.SetEntry(newEntry)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}
