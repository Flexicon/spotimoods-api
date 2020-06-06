package spotify

import (
	"fmt"
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
	return &Client{}
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
