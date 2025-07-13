package application

import (
	"encoding/json"
	"equestrian-events-api/internal/domain"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/petegabriel/hippobase"
	"net/http"
	"strconv"
	"time"
)

// ListCompetitions is a handler that returns a list of competitions
// from the in-memory cache or from the module
// if the cache is empty
// @Summary List competitions
// @Description List competitions
// @Produce JSON
// @Success 200 {array} CompetitionList
// @Router /competitions [get]
func (app *Application) ListCompetitions(_ *gin.Context) (CompetitionList, error) {

	events, err := hippobase.GetEvents()
	if err != nil {
		// TODO log the error
		return nil, New(http.StatusServiceUnavailable, "Unable to fetch events from Hippobase")
	}

	var competitions CompetitionList

	for _, parsed := range events {

		comp := &domain.Competition{
			ID:           parsed.Id,
			Name:         fmt.Sprintf("%s - %s", parsed.Location, parsed.Name),
			URL:          parsed.EventURL,
			Events:       make([]*domain.Event, 0),
			EntryListURL: parsed.EntryListURL,
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

// GetCompetitionByID is a handler that returns competition details
// from the in-memory cache or from the module based on the given ID.
// @Summary List competitions
// @Description List competitions
// @Produce JSON
// @Success 200 {object} Competition
// @Router /competitions/{id} [get]
func (app *Application) GetCompetitionByID(c *gin.Context) (Competition, error) {

	var comp *domain.Competition
	id := c.Param("id")

	err := app.InMemory.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return NotFound("Competition not found")
		}

		if item != nil {
			err := item.Value(func(val []byte) error {
				err := json.Unmarshal(val, &comp)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	entryList, err := hippobase.GetEntryLists(comp.EntryListURL)
	if err != nil {
		return nil, InternalError(err.Error())
	}

	if entryList != nil {
		for _, list := range entryList.Events {
			// info about the event
			event := &domain.Event{
				ID:       uuid.New().String(),
				Date:     list.CreatedAt,
				Name:     list.EventFullName,
				Nations:  list.TotalNations,
				Athletes: list.TotalAthletes,
				Horses:   list.TotalHorses,
				RidersAndHorses: make([]domain.Competitor, 0),
			}
			// info about each pair rider/horses
			// TODO maybe we can add the country info
			for _, competitors := range list.Competitors {
				for _, rider := range competitors.Pairs {
					competitor := domain.Competitor{
						Rider:  rider.Competitor,
						Horses: rider.Horses,
					}
					event.RidersAndHorses = append(event.RidersAndHorses, competitor)
				}
			}
			comp.Events = append(comp.Events, event)
		}
	}

	return comp, err
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

			idStr := strconv.Itoa(cpt.ID)
			newEntry := badger.NewEntry([]byte(idStr), mEvt).WithTTL(time.Hour)
			err = txn.SetEntry(newEntry)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}
