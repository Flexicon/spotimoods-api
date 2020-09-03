package internal

import "time"

// SpotifyToken for a particular user
type SpotifyToken struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Token     string `gorm:"unique;not null"`
	Refresh   string `gorm:"unique;not null"`
	UserID    uint
	User      User
}

// SpotifyProfile represents a user profile within Spotify
type SpotifyProfile struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	ID          string `json:"id"`
	URI         string `json:"uri"`

	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`

	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`

	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

// SpotifyTokenResponse retrieved from authorizing with the API
type SpotifyTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// SpotifyArtist response structure
//
// Docs: https://developer.spotify.com/documentation/web-api/reference/artists/
type SpotifyArtist struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Link         string         `json:"href"`
	Images       []SpotifyImage `json:"images"`
	Genres       []string       `json:"genres"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

// SpotifyImage structure
type SpotifyImage struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

// CreatePlaylistResponse from the spotify API
//
// Docs: https://developer.spotify.com/documentation/web-api/reference/playlists/create-playlist/
type CreatePlaylistResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SpotifyClient for all comunication with Spotify and it's API
type SpotifyClient interface {
	// GetAuthorizeURL prepares a url to begin the OAuth process with Spotify
	GetAuthorizeURL(state string) string
	// GetMyProfile fetches the user profile for the currently logged in user
	GetMyProfile(token *SpotifyToken) (*SpotifyProfile, error)
	// AuthorizeByCode with Spotify and return a token response
	AuthorizeByCode(code string) (*SpotifyTokenResponse, error)
	// CreatePlaylist makes a new playlist for the authed user and returns it's ID
	CreatePlaylist(token *SpotifyToken, name string) (string, error)
	// UpdatePlaylist edits an existing playlist for the authed user
	UpdatePlaylist(token *SpotifyToken, id, name string) error
	// DeletePlaylist for the authed user
	DeletePlaylist(token *SpotifyToken, id string) error
	// SearchForArtists by the given query
	SearchForArtists(token *SpotifyToken, query string) ([]*SpotifyArtist, error)
	// GetTopArtists for the user
	GetTopArtists(token *SpotifyToken) ([]*SpotifyArtist, error)
}
