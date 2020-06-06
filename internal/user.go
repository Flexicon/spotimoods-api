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
	FindByEmail(email string) (*User, error)
	Save(user *User) error
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
