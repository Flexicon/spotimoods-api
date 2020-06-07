package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/spf13/viper"
)

// Client for all comunication with Spotify and it's API
type Client struct {
	http internal.HTTPClient
}

// NewClient constructor
func NewClient(h internal.HTTPClient) *Client {
	return &Client{
		http: h,
	}
}

// GetAuthorizeURL prepares a url to begin the OAuth process with Spotify
func (c *Client) GetAuthorizeURL(state string) string {
	clientID := viper.GetString("spotify.client_id")
	apiDomain := viper.GetString("domains.api")

	q := url.Values{}
	q.Add("response_type", "code")
	q.Add("scope", "user-read-email user-top-read user-read-currently-playing user-read-recently-played")
	q.Add("client_id", clientID)
	q.Add("state", state)

	redirectURI := fmt.Sprintf("%s/callback", apiDomain)
	authorizeURL := fmt.Sprintf("https://accounts.spotify.com/authorize?redirect_uri=%s&%s", redirectURI, q.Encode())

	return authorizeURL
}

// GetMyProfile fetches the user profile for the currently logged in user
func (c *Client) GetMyProfile(token string) (*internal.SpotifyProfile, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user info: %v", err)
	}
	defer resp.Body.Close()

	var profile internal.SpotifyProfile
	json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, fmt.Errorf("Error parsing user info: %v", err)
	}

	return &profile, nil
}
