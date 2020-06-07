package internal

import "time"

// SpotifyToken for a particular user
type SpotifyToken struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Token     string `gorm:"unique;not null"`
	Refresh   string `gorm:"unique;not null"`
	UserID    int
	User      User
}

// SpotifyClient for all comunication with Spotify and it's API
type SpotifyClient interface {
	GetAuthorizeURL(state string) string
}
