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

// SpotifyClient for all comunication with Spotify and it's API
type SpotifyClient interface {
	// GetAuthorizeURL prepares a url to begin the OAuth process with Spotify
	GetAuthorizeURL(state string) string
	// GetMyProfile fetches the user profile for the currently logged in user
	GetMyProfile(token string) (*SpotifyProfile, error)
}
