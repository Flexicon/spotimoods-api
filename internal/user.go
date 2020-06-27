package internal

import (
	"github.com/jinzhu/gorm"
)

// User represents the applications user entity
type User struct {
	gorm.Model
	Email       string `gorm:"unique;not null"`
	DisplayName string
	Image       string
	SpotifyID   string `gorm:"not null"`
}

// UserRepository for interacting with user data
type UserRepository interface {
	// FindByEmail checks for an existing active user by a given email
	FindByEmail(email string) (*User, error)
	// FindTokenByUser attempts to retrieve a SpotifyToken for the given user ID
	FindTokenByUser(userID uint) (*SpotifyToken, error)
	// Save upserts the given user into the DB
	Save(user *User) error
	// SaveTokenForUser persists a new token or updets it for a given user
	SaveTokenForUser(user *User, token, refresh string) error
}

// UserService for performing all operations related to users
type UserService struct {
	r UserRepository
}

// NewUserService constructor
func NewUserService(r UserRepository) *UserService {
	return &UserService{
		r: r,
	}
}

// FindByEmail checks for an existing active user by a given email
func (s *UserService) FindByEmail(email string) (*User, error) {
	return s.r.FindByEmail(email)
}

// UpsertUser either updates or sets up a user with the given data and persists them
// User is identified by email
func (s *UserService) UpsertUser(id, name, email, image, token, refresh string) (*User, error) {
	user, err := s.FindByEmail(email)
	if err != nil && err != ErrNotFound {
		return nil, err
	}
	if user == nil {
		user = &User{}
	}

	user.DisplayName = name
	user.Email = email
	user.Image = image
	user.SpotifyID = id

	if err := s.r.Save(user); err != nil {
		return nil, err
	}

	return user, s.r.SaveTokenForUser(user, token, refresh)
}

// FindTokenForUser finds a stored spotify OAuth token for the given user
func (s *UserService) FindTokenForUser(userID uint) (*SpotifyToken, error) {
	return s.r.FindTokenByUser(userID)
}
