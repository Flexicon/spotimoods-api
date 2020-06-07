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
	g = g.Group("", authMiddlewareChain(Options{Services: h.services})...)
	g.GET("/me", h.Me())
}

func (h *userController) Me() echo.HandlerFunc {
	type errResponse struct {
		Msg string `json:"message"`
	}

	return func(c echo.Context) error {
		user := c.Get("user").(*internal.User)
		token, err := h.services.User().FindTokenForUser(user)
		if err != nil {
			log.Printf("Couldn't find token for the current user (ID: %d): %v", user.ID, err)
			return c.JSON(http.StatusInternalServerError, errResponse{Msg: "Couldn't find token for the current user"})
		}

		profile, err := h.services.Spotify().GetMyProfile(token.Token)
		if err != nil {
			log.Println("Failed to retrieve user profile:", err)
			return c.JSON(http.StatusInternalServerError, errResponse{Msg: "Failed to retrieve user profile"})
		}

		return c.JSON(http.StatusOK, profile)
	}
}