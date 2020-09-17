package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/flexicon/spotimoods-go/internal/api/model"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type moodController struct {
	services *internal.ServiceProvider
}

func newMood(services *internal.ServiceProvider) Controller {
	return &moodController{
		services: services,
	}
}

func (h *moodController) Routes(g *echo.Group) {
	g = g.Group("/moods")
	useAuthMiddleware(g, Options{Services: h.services})

	g.GET("", h.List())
	g.POST("", h.Create())
	g.GET("/:id", h.Show())
	g.PUT("/:id", h.Update())
	g.DELETE("/:id", h.Delete())
}

func (h *moodController) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*internal.User)
		moods, err := h.services.Mood().GetMoods(user)
		if err != nil {
			log.Printf("Failed to get moods for user (ID: %d): %v", user.ID, err)
			return c.JSON(http.StatusInternalServerError, ErrResponse{Msg: "Failed to get moods"})
		}

		return c.JSON(http.StatusOK, moods)
	}
}

func (h *moodController) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &model.MoodPayload{}
		if err := c.Bind(payload); err != nil {
			log.Printf("Failed to bind request body: %v", err)
			return c.NoContent(http.StatusBadRequest)
		}

		if err := payload.Validate(); err != nil {
			log.Printf("Payload did not pass validation: %+v", payload)
			log.Printf("Validation error: %v", err)
			return c.JSON(http.StatusBadRequest, ErrResponse{Msg: err.Error()})
		}

		user := c.Get("user").(*internal.User)
		mood, err := h.services.Mood().AddMood(payload.Name, payload.Color, user)
		if err != nil {
			log.Printf("Failed to add mood: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrResponse{Msg: "Failed to add mood"})
		}

		return c.JSON(http.StatusOK, mood)
	}
}

func (h *moodController) Show() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user.spotify_token").(*internal.SpotifyToken)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return notFound(c, "mood")
		}

		mood, err := h.services.Mood().FindForUser(uint(id), &token.User)
		if err != nil {
			return notFound(c, "mood")
		}

		// Populate artist data in mood tags
		artistIDs := make([]string, 0)
		for _, tag := range mood.Tags {
			artistIDs = append(artistIDs, tag.ArtistID)
		}

		artists, err := h.services.Spotify().GetArtistsByIDs(token, artistIDs)
		if err != nil {
			err := errors.Wrap(err, "failed to retrieve mood artist data")
			return c.JSON(http.StatusInternalServerError, ErrResponse{Msg: err.Error()})
		}

		for i, artist := range artists {
			mood.Tags[i].ArtistData = *artist
		}

		return c.JSON(http.StatusOK, mood)
	}
}

func (h *moodController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*internal.User)
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return notFound(c, "mood")
		}

		payload := &model.MoodChanges{}
		if err := c.Bind(payload); err != nil {
			log.Printf("Failed to bind request body: %v", err)
			return c.NoContent(http.StatusBadRequest)
		}

		if err := payload.Validate(); err != nil {
			log.Printf("Payload did not pass validation: %+v", payload)
			log.Printf("Validation error: %v", err)
			return c.JSON(http.StatusBadRequest, ErrResponse{Msg: err.Error()})
		}

		changes := internal.Mood{
			Name:  payload.Name,
			Color: payload.Color,
		}

		mood, err := h.services.Mood().UpdateMoodForUser(uint(id), changes, user)
		if err != nil {
			if err == internal.ErrNotFound {
				return notFound(c, "mood")
			}
			log.Printf("Failed to update mood: %v", err)
			log.Printf("Payload: %+v", payload)
			return c.JSON(http.StatusInternalServerError, ErrResponse{Msg: "failed to update mood"})
		}

		return c.JSON(http.StatusOK, mood)
	}
}

func (h *moodController) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*internal.User)
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("failed to convert id to int: %v", err)
			return notFound(c, "mood")
		}

		err = h.services.Mood().DeleteForUser(uint(id), user)
		if err != nil {
			if err == internal.ErrNotFound {
				return notFound(c, "mood")
			}
			log.Printf("Failed to delete mood: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrResponse{Msg: "failed to delete mood"})
		}

		return c.NoContent(http.StatusOK)
	}
}
