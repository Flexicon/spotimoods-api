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
}

// UserRepository for interacting with user data
type UserRepository interface {
	// FindByEmail checks for an existing active user by a given email
	FindByEmail(email string) (*User, error)
	// FindTokenByUser attempts to retrieve a SpotifyToken for the given user
	FindTokenByUser(user *User) (*SpotifyToken, error)
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
func (s *UserService) UpsertUser(name, email, image, token, refresh string) error {
	user, err := s.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		user = &User{
			DisplayName: name,
			Email:       email,
			Image:       image,
		}
	}

	if err := s.r.Save(user); err != nil {
		return err
	}

	return s.r.SaveTokenForUser(user, token, refresh)
}

// FindTokenForUser finds a stored spotify OAuth token for the given user
func (s *UserService) FindTokenForUser(user *User) (*SpotifyToken, error) {
	return s.r.FindTokenByUser(user)
}
