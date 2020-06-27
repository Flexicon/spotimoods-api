package internal

import (
	"time"
)

// Mood represents a Mood entity
type Mood struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name       string    `gorm:"not null" json:"name"`
	Color      string    `gorm:"not null" json:"color"`
	PlaylistID string    `json:"playlist_id"`
	UserID     uint      `json:"-"`
	User       User      `json:"-"`
}

// MoodRepository for interacting with mood data
type MoodRepository interface {
	// Find mood by ID and User
	Find(id uint) (*Mood, error)
	// FindByIDAndUser if it exists
	FindByIDAndUser(id uint, user *User) (*Mood, error)
	// Remove mood by ID
	Remove(id uint) error
	// FindByUser all moods for a given user
	FindByUser(user *User) ([]*Mood, error)
	// Save upserts the given mood into the DB
	Save(mood *Mood) error
	// Update persists the given changeset to the given mood
	Update(mood *Mood, changes Mood) error
}

// MoodService for performing all operations related to moods
type MoodService struct {
	r MoodRepository
}

// NewMoodService constructor
func NewMoodService(r MoodRepository) *MoodService {
	return &MoodService{
		r: r,
	}
}

// AddMood for the given user
func (s *MoodService) AddMood(name, color string, user *User) (*Mood, error) {
	mood := &Mood{
		Name:  name,
		Color: color,
		User:  *user,
	}
	// TODO: Either create a playlist or kickoff a background worker to do it

	return mood, s.r.Save(mood)
}

// UpdateMoodForUser for a given change set
func (s *MoodService) UpdateMoodForUser(id uint, changes Mood, user *User) (*Mood, error) {
	mood, err := s.FindForUser(id, user)
	if err != nil {
		return nil, err
	}

	if err := s.r.Update(mood, changes); err != nil {
		return nil, err
	}
	// TODO: Either update the playlist name or kickoff a background worker to do it

	return mood, nil
}

// GetMoods finds all moods for a given user
func (s *MoodService) GetMoods(user *User) ([]*Mood, error) {
	return s.r.FindByUser(user)
}

// FindForUser finds a mood by the given ID and user
func (s *MoodService) FindForUser(id uint, user *User) (*Mood, error) {
	mood, err := s.r.FindByIDAndUser(id, user)
	if err != nil {
		return nil, err
	}
	return mood, nil
}

// DeleteForUser removes the stored mood by the given ID and user
func (s *MoodService) DeleteForUser(id uint, user *User) error {
	_, err := s.FindForUser(id, user)
	if err != nil {
		return err
	}
	// TODO: Either delete the related playlist or kickoff a background worker to do it

	return s.r.Remove(id)
}
