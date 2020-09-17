package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Client for all comunication with Spotify and it's API
type Client struct {
	http  internal.HTTPClient
	repos internal.RepositoryProvider
	cache internal.Cache
}

// NewClient constructor
func NewClient(h internal.HTTPClient, repos internal.RepositoryProvider, cache internal.Cache) *Client {
	return &Client{
		http:  h,
		repos: repos,
		cache: cache,
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
	cacheConfig := &internal.CacheItem{
		Key: fmt.Sprintf("GetMyProfile-%d", token.UserID),
		TTL: time.Minute,
	}

	body, err := c.fetchWithCache(req, token, cacheConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user info")
	}

	var profile internal.SpotifyProfile
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&profile); err != nil {
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
	payload, err := json.Marshal(PlaylistPayload{Name: name})
	if err != nil {
		return "", fmt.Errorf("failed to prepare payload: %v", err)
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", token.User.SpotifyID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to prepare request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.do(req, token)
	if err != nil {
		return "", fmt.Errorf("request failed when creating playlist: %v", err)
	}

	var playlist internal.CreatePlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return "", fmt.Errorf("Error parsing playlist response: %v", err)
	}
	if playlist.ID == "" {
		return "", fmt.Errorf("failed to create playlist id empty")
	}

	return playlist.ID, nil
}

// UpdatePlaylist edits an existing playlist for the authed user
func (c *Client) UpdatePlaylist(token *internal.SpotifyToken, id, name string) error {
	payload, err := json.Marshal(PlaylistPayload{Name: name})
	if err != nil {
		return fmt.Errorf("failed to prepare payload: %v", err)
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", id)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to prepare request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if _, err := c.do(req, token); err != nil {
		return fmt.Errorf("request failed when updating playlist: %v", err)
	}

	return nil
}

// DeletePlaylist really unfollows a given playlist ID, since spotify doesn't actually offer any way to delete a playlist
func (c *Client) DeletePlaylist(token *internal.SpotifyToken, id string) error {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/followers", id), nil)

	if _, err := c.do(req, token); err != nil {
		return fmt.Errorf("request failed when deleting playlist: %v", err)
	}

	return nil
}

// SearchForArtists by the given query
func (c *Client) SearchForArtists(token *internal.SpotifyToken, query string) ([]*internal.SpotifyArtist, error) {
	searchURL, _ := url.Parse("https://api.spotify.com/v1/search")
	q := url.Values{}
	q.Add("q", query)
	q.Add("type", "artist")
	searchURL.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, searchURL.String(), nil)
	cacheItem := &internal.CacheItem{
		Key: fmt.Sprintf("SearchForArtists-user-%d-%s", token.UserID, searchURL.RawQuery),
		TTL: time.Minute,
	}

	body, err := c.fetchWithCache(req, token, cacheItem)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve search results")
	}

	var searchResponse struct {
		Artists struct {
			Items []*internal.SpotifyArtist `json:"items"`
		} `json:"artists"`
	}
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("error parsing artists response: %v", err)
	}

	return searchResponse.Artists.Items, nil
}

// GetTopArtists for the user
func (c *Client) GetTopArtists(token *internal.SpotifyToken) ([]*internal.SpotifyArtist, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/top/artists", nil)
	cacheItem := &internal.CacheItem{
		Key: fmt.Sprintf("GetTopArtists-user-%d", token.UserID),
		TTL: time.Minute,
	}

	body, err := c.fetchWithCache(req, token, cacheItem)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve top artists")
	}

	var topResponse struct {
		Items []*internal.SpotifyArtist `json:"items"`
	}
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&topResponse); err != nil {
		return nil, fmt.Errorf("error parsing top artists response: %v", err)
	}

	return topResponse.Items, nil
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
