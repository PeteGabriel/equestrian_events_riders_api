package main

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	riders "github.com/petegabriel/equestrian_events_riders_list"
	"net/http"
)

func (a *Application) CheckCacheForEntryLists(c *gin.Context) {
	// check if data already exists
	// if so return it
	// if not, call the next handler
	if a.InMemory == nil {
		c.Next()
		return
	}

	var competitions []Competition

	err := a.InMemory.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			var event riders.EquestrianCompetition
			err := item.Value(func(val []byte) error {
				err := json.Unmarshal(val, &event)
				if err != nil {
					return err
				}
				c := Competition{Name: event.MainTitle, ID: event.MainTitle}
				c.Events = make([]*Event, 0)
				for _, e := range event.Events {
					c.Events = append(c.Events, &Event{
						Date:     e.CreatedAt,
						Name:     e.EventFullName,
						Nations:  e.TotalNations,
						Athletes: e.TotalAthletes,
						Horses:   e.TotalHorses,
					})
				}
				competitions = append(competitions, c)
				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		c.Abort()
		return
	}

	// if we found something in cache, return it
	if len(competitions) > 0 {
		c.JSON(http.StatusNotModified, gin.H{
			"message": competitions,
		})
		c.Abort()
		return
	}

	c.Next()
}
