package application

import (
	"equestrian-events-api/internal/domain"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"log/slog"
	"net/http"
)

type HTTPHandlerWithErr[T any] func(c *gin.Context) (T, error)

type CompetitionList []*domain.Competition

type Competition *domain.Competition

// HandleCompetitions is a Gin handler that processes the request for competitions.
func (app *Application) HandleCompetitions(handler HTTPHandlerWithErr[CompetitionList]) gin.HandlerFunc {
	return func(c *gin.Context) {

		var competitions []*domain.Competition
		competitions, err := handler(c)

		c.Writer.Header().Set("Content-Type", jsonapi.MediaType)

		if err != nil {
			handleHTTPError(c, err)
			return
		}

		if err := jsonapi.MarshalPayload(c.Writer, competitions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
		}

	}
}

func (app *Application) HandleCompetitionByID(handler HTTPHandlerWithErr[Competition]) gin.HandlerFunc {
	return func(c *gin.Context) {

		var competition Competition
		competition, err := handler(c)

		c.Writer.Header().Set("Content-Type", jsonapi.MediaType)

		if err != nil {
			handleHTTPError(c, err)
			return
		}

		if err := jsonapi.MarshalPayload(c.Writer, competition); err != nil {
			slog.Error("Failed to marshal competition inside handler: %v", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
		}

	}
}

func handleHTTPError(c *gin.Context, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		slog.Debug("http error", "code", httpErr.Code, slog.String("error", err.Error()))
		//TODO observe and act accordingly
		c.JSON(httpErr.Code, gin.H{
			"message": err.Error(),
			"code":    httpErr.Code,
			"err":     err.Error(),
		})
	} else {
		// Default to 500
		slog.Error("internal server error", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
			"code":    httpErr.Code,
			"err":     err.Error(),
		})

	}
}
