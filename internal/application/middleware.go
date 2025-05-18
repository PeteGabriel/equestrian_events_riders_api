package application

import (
	"encoding/json"
	"equestrian-events-api/internal/domain"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"net/http"
)

// CheckCacheForEntryLists checks if the data is already in the cache.
// If it is, it returns the data from the cache.
// If not, it calls the next handler in the chain.
func (app *Application) CheckCacheForEntryLists(c *gin.Context) {
	if app.InMemory == nil {
		c.Next()
		return
	}

	var competitions []*domain.Competition

	err := app.InMemory.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			var event *domain.Competition
			err := item.Value(func(val []byte) error {
				err := json.Unmarshal(val, &event)
				if err != nil {
					return err
				}
				c := &domain.Competition{Name: event.Name, ID: event.ID}
				c.Events = make([]*domain.Event, 0)
				for _, e := range event.Events {
					c.Events = append(c.Events, &domain.Event{
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
