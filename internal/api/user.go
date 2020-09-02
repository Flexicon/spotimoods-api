package api

import (
	"log"
	"net/http"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/labstack/echo/v4"
)

type userController struct {
	services *internal.ServiceProvider
}

func newUser(services *internal.ServiceProvider) Controller {
	return &userController{
		services: services,
	}
}

func (h *userController) Routes(g *echo.Group) {
	g = g.Group("/users")
	useAuthMiddleware(g, Options{Services: h.services})

	g.GET("/me", h.Me())
}

func (h *userController) Me() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user.spotify_token").(*internal.SpotifyToken)
		profile, err := h.services.Spotify().GetMyProfile(token)
		if err != nil {
			log.Println("Failed to retrieve user profile:", err)
			return c.JSON(http.StatusInternalServerError, ErrResponse{Msg: "Failed to retrieve user profile"})
		}

		return c.JSON(http.StatusOK, profile)
	}
}
