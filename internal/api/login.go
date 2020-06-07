package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

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
	// TODO: replace with calls to SpotifyClient
	client := &http.Client{Timeout: 5 * time.Second}

	return func(c echo.Context) error {
		q := c.QueryParams()
		code := q.Get("code")
		state := q.Get("state")
		clientID := viper.GetString("spotify.client_id")
		clientSecret := viper.GetString("spotify.client_secret")
		apiDomain := viper.GetString("domains.api")

		// TODO: comapre with state stored in cookie/cache
		if code == "" || state == "" || state != "123" {
			return c.String(http.StatusBadRequest, "State mismatch")
		}

		form := url.Values{}
		form.Set("code", code)
		form.Set("grant_type", "authorization_code")
		form.Set("redirect_uri", fmt.Sprintf("%s/callback", apiDomain))

		tokenURL := "https://accounts.spotify.com/api/token"
		req, _ := http.NewRequest(http.MethodPost, tokenURL, bytes.NewBuffer([]byte(form.Encode())))

		authorizationToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
		req.Header.Set("Authorization", "Basic "+authorizationToken)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Error authorizing token:", err.Error()))
		}
		defer resp.Body.Close()

		type spotifyTokenResponse struct {
			AccessToken      string `json:"access_token"`
			RefreshToken     string `json:"refresh_token"`
			ExpiresIn        int    `json:"expires_in"`
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		var tokenResp spotifyTokenResponse
		json.NewDecoder(resp.Body).Decode(&tokenResp)

		if tokenResp.Error != "" {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Failed to login:", tokenResp.ErrorDescription))
		}

		profile, err := h.services.Spotify().GetMyProfile(tokenResp.AccessToken)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Error parsing user info:", err))
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

		err = h.services.User().UpsertUser(profile.DisplayName, profile.Email, image, tokenResp.AccessToken, tokenResp.RefreshToken)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintln("Failed to register user:", err))
		}

		// TODO: redirect to app
		return c.Redirect(http.StatusFound, fmt.Sprintf("%s/api/ping?msg=%s", apiDomain, signedToken))
	}
}
