package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/labstack/echo/v4"
)

type spotifyController struct {
	services *internal.ServiceProvider
}

func newSpotify(services *internal.ServiceProvider) Controller {
	return &spotifyController{
		services: services,
	}
}

func (h *spotifyController) Routes(g *echo.Group) {
	artists := g.Group("/artists")
	useAuthMiddleware(artists, Options{Services: h.services})

	artists.GET("/search", h.ArtistsSearch())
	artists.GET("/top", h.TopArtists())
}

func (h *spotifyController) ArtistsSearch() echo.HandlerFunc {
	return func(c echo.Context) error {
		q := c.QueryParam("query")
		if q == "" {
			return c.JSON(http.StatusBadRequest, ErrResponse{Msg: "search query is required"})
		}

		token := c.Get("user.spotify_token").(*internal.SpotifyToken)
		artists, err := h.services.Spotify().SearchForArtists(token, q)
		if err != nil {
			errMsg := fmt.Sprintf("failed to search for artists: %v", err)
			log.Printf(errMsg)
			return c.JSON(http.StatusInternalServerError, ErrResponse{Msg: errMsg})
		}

		return c.JSON(http.StatusOK, artists)
	}
}

func (h *spotifyController) TopArtists() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user.spotify_token").(*internal.SpotifyToken)

		artists, err := h.services.Spotify().GetTopArtists(token)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get top artists: %v", err)
			log.Printf(errMsg)
			return c.JSON(http.StatusInternalServerError, ErrResponse{Msg: errMsg})
		}

		return c.JSON(http.StatusOK, artists)
	}
}
