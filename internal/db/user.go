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
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// Save upserts the given user into the DB
func (r *UserRepository) Save(user *internal.User) error {
	if r.db.NewRecord(user) {
		return r.db.Create(&user).Error
	}
	return r.db.Save(&user).Error
}

// SaveTokenForUser persists a new token for a given user
func (r *UserRepository) SaveTokenForUser(user *internal.User, token, refresh string) error {
	spotToken := &internal.SpotifyToken{
		Token:   token,
		Refresh: refresh,
		User:    *user,
	}

	return r.db.Create(&spotToken).Error
}
