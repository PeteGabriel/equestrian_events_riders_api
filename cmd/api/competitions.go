package main

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	riders "github.com/petegabriel/hippobase"
	"time"
)

// ListCompetitions is a handler that returns a list of competitions
// from the in-memory cache or from the module
// if the cache is empty
// @Summary List competitions
// @Description List competitions
// @Produce json
// @Success 200 {array} []EquineCompetition
// @Router /competitions [get]
func (a *Application) ListCompetitions() ([]*Competition, error) {

	parsedComps := riders.Parse()

	var competitions []*Competition

	for _, parsed := range parsedComps {
		comp := &Competition{
			ID:     uuid.New().String(),
			Name:   parsed.MainTitle,
			Events: make([]*Event, 0),
		}
		for _, evt := range parsed.Events {
			comp.Events = append(comp.Events, &Event{
				ID:          uuid.New().String(),
				Date:        evt.CreatedAt,
				Name:        evt.EventFullName,
				Nations:     evt.TotalNations,
				Athletes:    evt.TotalAthletes,
				Horses:      evt.TotalHorses,
				Competitors: evt.Competitors,
			})
		}
		competitions = append(competitions, comp)
	}

	//TODO explore the possibility of doing this in a different goroutine
	if err := a.cacheEvents(competitions); err != nil {
		return nil, InternalError(err.Error())
	}

	return competitions, nil
}

func (a *Application) cacheEvents(cpts []*Competition) error {
	if a.InMemory == nil {
		return nil
	}
	for _, cpt := range cpts {
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
