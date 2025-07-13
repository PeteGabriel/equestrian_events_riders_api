package application

import (
	"encoding/json"
	"equestrian-events-api/internal/domain"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
)

// CheckCompetitionsInCache retrieves competitions from the cache if available,
// otherwise proceeds with the next handler.
func (app *Application) CheckCompetitionsInCache(c *gin.Context) {
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
					slog.Error("Error unmarshalling event", slog.String("error", err.Error()))
					return err
				}
				competitions = append(competitions, event)
				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		c.Abort()
		return
	}

	if len(competitions) > 0 {
		slog.Info("total competitions found in cache", slog.String("total_competitions", strconv.Itoa(len(competitions))))
		c.Writer.Header().Set("Content-Type", jsonapi.MediaType)
		c.Writer.WriteHeader(http.StatusOK)

		if err = jsonapi.MarshalPayload(c.Writer, competitions); err != nil {
			slog.Error("failed to marshal competitions payload", slog.String("error", err.Error()))
			c.JSON(http.StatusOK, gin.H{
				"message": "error retrieving data from cache",
			})
		}

		c.Abort()
		return
	}

	c.Next()
}
