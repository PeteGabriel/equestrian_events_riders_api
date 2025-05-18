package main

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
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

	var competitions []*Competition

	err := a.InMemory.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			var event *Competition
			err := item.Value(func(val []byte) error {
				err := json.Unmarshal(val, &event)
				if err != nil {
					return err
				}
				c := &Competition{Name: event.Name, ID: event.ID}
				c.Events = make([]*Event, 0)
				for _, e := range event.Events {
					c.Events = append(c.Events, &Event{
						ID:          e.ID,
						Date:        e.Date,
						Name:        e.Name,
						Nations:     e.Nations,
						Athletes:    e.Athletes,
						Horses:      e.Horses,
						Competitors: e.Competitors,
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
		c.Writer.Header().Set("Content-Type", jsonapi.MediaType)
		c.Writer.WriteHeader(http.StatusOK)

		if err = jsonapi.MarshalPayload(c.Writer, competitions); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "error retrieving data from cache",
			})
		}

		c.Abort()
		return
	}

	c.Next()
}
