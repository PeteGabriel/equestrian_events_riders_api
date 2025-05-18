package application

import (
	"equestrian-events-api/internal/domain"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"log/slog"
	"net/http"
)

type HTTPHandlerWithErr[T any] func() (T, error)

type CompetitionList []*domain.Competition

func (app *Application) HandleCompetitions(handler HTTPHandlerWithErr[CompetitionList]) gin.HandlerFunc {
	return func(c *gin.Context) {

		var competitions []*domain.Competition
		competitions, err := handler()

		if err != nil {
			// Check if it's an HTTPError
			var httpErr *HTTPError
			if errors.As(err, &httpErr) {
				slog.Debug("http error", "code", httpErr.Code, "err", err.Error())
				//TODO observe and act accordingly
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal server error",
				})
			} else {
				// Default to 500
				slog.Error("internal server error", "err", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal server error",
				})
			}
		}

		c.Writer.Header().Set("Content-Type", jsonapi.MediaType)
		c.Writer.WriteHeader(http.StatusOK)

		if err := jsonapi.MarshalPayload(c.Writer, competitions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
		}
	}
}
