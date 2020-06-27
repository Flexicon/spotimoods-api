package db

import (
	"github.com/flexicon/spotimoods-go/internal"
	"github.com/jinzhu/gorm"
)

// UserRepository for interacting with user data in the DB
type UserRepository struct {
	db *gorm.DB
}

// FindByEmail checks for an existing active user by a given email
func (r *UserRepository) FindByEmail(email string) (*internal.User, error) {
	var user internal.User
	query := r.db.Where("email = ?", email).First(&user)
	if query.RecordNotFound() {
		return nil, internal.ErrNotFound
	}

	return &user, query.Error
}

// FindTokenByUser attempts to retrieve a SpotifyToken for the given user ID
func (r *UserRepository) FindTokenByUser(userID uint) (*internal.SpotifyToken, error) {
	var token internal.SpotifyToken
	query := r.db.Preload("User").Where("user_id = ?", userID).First(&token)
	if query.RecordNotFound() {
		return nil, internal.ErrNotFound
	}

	return &token, query.Error
}

// Save upserts the given user into the DB
func (r *UserRepository) Save(user *internal.User) error {
	if r.db.NewRecord(user) {
		return r.db.Create(&user).Error
	}
	return r.db.Save(&user).Error
}

// SaveTokenForUser persists a new token or updates it for a given user
func (r *UserRepository) SaveTokenForUser(user *internal.User, token, refresh string) error {
	var spotToken internal.SpotifyToken
	r.db.Where("user_id = ?", user.ID).First(&spotToken)

	spotToken.Token = token
	spotToken.UserID = user.ID
	if refresh != "" {
		spotToken.Refresh = refresh
	}

	if r.db.NewRecord(spotToken) {
		return r.db.Create(&spotToken).Error
	}
	return r.db.Save(&spotToken).Error
}
