package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/spf13/viper"
)

// Client for all comunication with Spotify and it's API
type Client struct {
	http  internal.HTTPClient
	repos internal.RepositoryProvider
}

// NewClient constructor
func NewClient(h internal.HTTPClient, repos internal.RepositoryProvider) *Client {
	return &Client{
		http:  h,
		repos: repos,
	}
}

// GetAuthorizeURL prepares a url to begin the OAuth process with Spotify
func (c *Client) GetAuthorizeURL(state string) string {
	clientID := viper.GetString("spotify.client_id")
	scope := viper.GetString("spotify.scope")
	apiDomain := viper.GetString("domains.api")

	q := url.Values{}
	q.Add("response_type", "code")
	q.Add("scope", scope)
	q.Add("client_id", clientID)
	q.Add("state", state)
	q.Add("redirect_uri", fmt.Sprintf("%s/callback", apiDomain))

	url, _ := url.Parse("https://accounts.spotify.com/authorize")
	url.RawQuery = q.Encode()

	return url.String()
}

// GetMyProfile fetches the user profile for the currently logged in user
func (c *Client) GetMyProfile(token *internal.SpotifyToken) (*internal.SpotifyProfile, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+token.Token)

	resp, err := c.do(req, token)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user info: %v", err)
	}
	defer resp.Body.Close()

	var profile internal.SpotifyProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("Error parsing user info: %v", err)
	}

	return &profile, nil
}

// AuthorizeByCode with Spotify and return a token response
func (c *Client) AuthorizeByCode(code string) (*internal.SpotifyTokenResponse, error) {
	return c.Authorize(code, "code", "authorization_code")
}

// Authorize with Spotify and return a token response
func (c *Client) Authorize(grant, grantName, grantType string) (*internal.SpotifyTokenResponse, error) {
	clientID := viper.GetString("spotify.client_id")
	clientSecret := viper.GetString("spotify.client_secret")
	apiDomain := viper.GetString("domains.api")

	form := url.Values{}
	form.Set(grantName, grant)
	form.Set("grant_type", grantType)
	form.Set("redirect_uri", fmt.Sprintf("%s/callback", apiDomain))

	tokenURL := "https://accounts.spotify.com/api/token"
	req, err := http.NewRequest(http.MethodPost, tokenURL, bytes.NewBuffer([]byte(form.Encode())))

	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %v", err)
	}

	authorizationToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Set("Authorization", "Basic "+authorizationToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to authorize with spotify: %v", err)
	}
	defer resp.Body.Close()

	var token internal.SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to parse spotify token response: %v", err)
	}

	return &token, nil
}

// CreatePlaylist makes a new playlist for the authed user and returns it's ID
func (c *Client) CreatePlaylist(token *internal.SpotifyToken, name string) (string, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/users/1172023744/playlists", nil)

	resp, err := c.do(req, token)
	if err != nil {
		return "", fmt.Errorf("Error fetching user info: %v", err)
	}
	defer resp.Body.Close()

	var playlist internal.CreatePlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return "", fmt.Errorf("Error parsing playlist response: %v", err)
	}

	return playlist.ID, nil
}

// Refresh the given token with spotify
func (c *Client) Refresh(token *internal.SpotifyToken) error {
	st, err := c.Authorize(token.Refresh, "refresh_token", "refresh_token")
	if err != nil {
		return err
	}

	if err := c.repos.User().SaveTokenForUser(&token.User, st.AccessToken, st.RefreshToken); err != nil {
		return err
	}

	token.Token = st.AccessToken
	if st.RefreshToken != "" {
		token.Refresh = st.RefreshToken
	}

	return nil
}

// do performs an action against the spotify API and on authorization failure attempts to refresh the given token and try again
func (c *Client) do(req *http.Request, token *internal.SpotifyToken) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+token.Token)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusUnauthorized || token.Refresh == "" {
		if resp.StatusCode >= 400 {
			return nil, httpStatusErr(resp)
		}
		return resp, nil
	}

	if err := c.Refresh(token); err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.Token)
	resp, err = c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, httpStatusErr(resp)
	}

	return resp, nil
}

func httpStatusErr(resp *http.Response) error {
	body, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("Http status %d: %s", resp.StatusCode, body)
}
