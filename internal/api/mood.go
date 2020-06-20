package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/labstack/echo/v4"
)

type moodController struct {
	services *internal.ServiceProvider
}

// MoodPayload for creating a new Mood
type MoodPayload struct {
	Name  string `json:"name"`
	Color string `json:"color"`
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
	g.DELETE("/:id", h.Delete())
}

func (h *moodController) List() echo.HandlerFunc {
	type errResponse struct {
		Msg string `json:"message"`
	}

	return func(c echo.Context) error {
		user := c.Get("user").(*internal.User)
		moods, err := h.services.Mood().GetMoods(user)
		if err != nil {
			log.Printf("Failed to get moods for user (ID: %d): %v", user.ID, err)
			return c.JSON(http.StatusInternalServerError, errResponse{Msg: "Failed to get moods"})
		}

		return c.JSON(http.StatusOK, moods)
	}
}

func (h *moodController) Create() echo.HandlerFunc {
	type errResponse struct {
		Msg string `json:"message"`
	}

	return func(c echo.Context) error {
		payload := &MoodPayload{}
		if err := c.Bind(payload); err != nil {
			log.Printf("Failed to bind request body: %v", err)
			return c.NoContent(http.StatusBadRequest)
		}

		if payload.Name == "" || payload.Color == "" {
			log.Printf("Payload did not pass validation: %+v", payload)
			return c.JSON(http.StatusBadRequest, errResponse{Msg: "name and color is required"})
		}

		user := c.Get("user").(*internal.User)
		mood, err := h.services.Mood().AddMood(payload.Name, payload.Color, user)
		if err != nil {
			log.Printf("Failed to add mood: %v", err)
			return c.JSON(http.StatusInternalServerError, errResponse{Msg: "Failed to add mood"})
		}

		return c.JSON(http.StatusOK, mood)
	}
}

func (h *moodController) Delete() echo.HandlerFunc {
	type errResponse struct {
		Msg string `json:"message"`
	}

	return func(c echo.Context) error {
		user := c.Get("user").(*internal.User)
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, errResponse{Msg: fmt.Sprintf("Invalid mood id '%s'", c.Param("id"))})
		}

		err = h.services.Mood().DeleteForUser(uint(id), user)
		if err != nil {
			if err == internal.ErrNotFound {
				return c.JSON(http.StatusNotFound, errResponse{Msg: "mood not found"})
			}
			log.Printf("Failed to delete mood: %v", err)
			return c.JSON(http.StatusInternalServerError, errResponse{Msg: "failed to delete mood"})
		}

		return c.NoContent(http.StatusOK)
	}
}
