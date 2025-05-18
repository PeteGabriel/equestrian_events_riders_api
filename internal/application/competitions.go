package application

import (
	"encoding/json"
	"equestrian-events-api/internal/domain"
	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/petegabriel/hippobase"
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
func (app *Application) ListCompetitions() (CompetitionList, error) {

	parsedComps := hippobase.Parse()

	var competitions CompetitionList

	for _, parsed := range parsedComps {
		comp := &domain.Competition{
			ID:     uuid.New().String(),
			Name:   parsed.MainTitle,
			Events: make([]*domain.Event, 0),
		}
		for _, evt := range parsed.Events {
			comp.Events = append(comp.Events, &domain.Event{
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
	if err := cacheEvents(competitions, app); err != nil {
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
