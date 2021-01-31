package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/flexicon/spotimoods-go/internal/api/auth"
	"github.com/google/uuid"
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
		// Generate state key
		state, err := uuid.NewRandom()
		if err != nil {
			log.Printf("Failed to generate UUID for state param: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to generate state")
		}

		// Persist state as HttpOnly cookie
		cookieExpiry := time.Now().Add(5 * time.Minute)
		c.SetCookie(&http.Cookie{
			Name:     "state",
			Value:    state.String(),
			HttpOnly: true,
			Expires:  cookieExpiry,
		})

		// Redirect to spotify auth
		return c.Redirect(http.StatusFound, h.services.Spotify().GetAuthorizeURL(state.String()))
	}
}

func (h *loginController) loginCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		q := c.QueryParams()
		code := q.Get("code")
		state := q.Get("state")

		// State validation
		stateCookie, err := c.Cookie("state")
		if err != nil {
			log.Printf("Failed to retrieve state cookie: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to validate state")
		}
		if code == "" || state == "" || state != stateCookie.Value {
			return c.String(http.StatusBadRequest, "State mismatch")
		}

		// Spotify Auth
		token, err := h.services.Spotify().AuthorizeByCode(code)
		if err != nil {
			log.Printf("failed to authorize with spotify: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to authorize with spotify")
		}
		if token.Error != "" {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Failed to login:", token.ErrorDescription))
		}

		// Fetch spotify profile
		profile, err := h.services.Spotify().GetMyProfile(&internal.SpotifyToken{Token: token.AccessToken})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Generate app JWT
		signedToken, err := auth.GenerateToken(auth.TokenOptions{
			DisplayName: profile.DisplayName,
			Email:       profile.Email,
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Error signing token:", err))
		}

		// Prepare and persist user to system
		var image string
		if len(profile.Images) > 0 {
			image = profile.Images[0].URL
		}

		_, err = h.services.User().UpsertUser(profile.ID, profile.DisplayName, profile.Email, image, token.AccessToken, token.RefreshToken)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Failed to register user:", err))
		}

		// Redirect back to web app with token
		return c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth?token=%s", viper.GetString("domains.front"), signedToken))
	}
}
