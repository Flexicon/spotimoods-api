package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type loginController struct {
	services *internal.ServiceProvider
}

func newLogin(services *internal.ServiceProvider) Controller {
	return &loginController{
		services: services,
	}
}

func (h *loginController) Routes(g *echo.Group) {
	g.GET("/login", h.login())
	g.GET("/callback", h.loginCallback())
}

func (h *loginController) login() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: generate and store state in cookie/cache
		return c.Redirect(http.StatusFound, h.services.Spotify().GetAuthorizeURL("123"))
	}
}

func (h *loginController) loginCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		q := c.QueryParams()
		code := q.Get("code")
		state := q.Get("state")

		// TODO: comapre with state stored in cookie/cache
		if code == "" || state == "" || state != "123" {
			return c.String(http.StatusBadRequest, "State mismatch")
		}

		token, err := h.services.Spotify().AuthorizeByCode(code)
		if err != nil {
			log.Printf("failed to authorize with spotify: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to authorize with spotify")
		}
		if token.Error != "" {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Failed to login:", token.ErrorDescription))
		}

		profile, err := h.services.Spotify().GetMyProfile(&internal.SpotifyToken{Token: token.AccessToken})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		signedToken, err := generateToken(TokenOptions{
			DisplayName: profile.DisplayName,
			Email:       profile.Email,
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Error signing token:", err))
		}

		var image string
		if len(profile.Images) > 0 {
			image = profile.Images[0].URL
		}

		_, err = h.services.User().UpsertUser(profile.ID, profile.DisplayName, profile.Email, image, token.AccessToken, token.RefreshToken)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Failed to register user:", err))
		}

		// TODO: redirect to app
		return c.Redirect(http.StatusFound, fmt.Sprintf("%s/api/ping?msg=%s", viper.GetString("domains.api"), signedToken))
	}
}
