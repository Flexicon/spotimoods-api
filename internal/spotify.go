package internal

import "github.com/jinzhu/gorm"

// SpotifyToken for a particular user
type SpotifyToken struct {
	gorm.Model
	Token   string `gorm:"unique;not null"`
	Refresh string `gorm:"unique;not null"`
	User    User
}

// SpotifyClient for all comunication with Spotify and it's API
type SpotifyClient interface {
	GetAuthorizeURL(state string) string
}
