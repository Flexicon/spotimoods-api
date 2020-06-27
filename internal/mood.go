package internal

import (
	"encoding/json"
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

// MarshalJSON for api responses
func (m *Mood) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Mood
		HasPlaylist bool `json:"has_playlist"`
	}{
		*m,
		m.PlaylistID != "",
	})
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
	r       MoodRepository
	q       QueueService
	spotify SpotifyClient
}

// NewMoodService constructor
func NewMoodService(r MoodRepository, q QueueService, s SpotifyClient) *MoodService {
	return &MoodService{
		r:       r,
		q:       q,
		spotify: s,
	}
}

// AddMood for the given user
func (s *MoodService) AddMood(name, color string, user *User) (*Mood, error) {
	mood := &Mood{
		Name:  name,
		Color: color,
		User:  *user,
	}
	if err := s.r.Save(mood); err != nil {
		return nil, err
	}

	// Add task to create playlist in spotify
	if err := s.q.AddPlaylist(mood); err != nil {
		return nil, err
	}

	return mood, nil
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

	// Add task to update playlist in spotify
	if err := s.q.UpdatePlaylist(mood); err != nil {
		return nil, err
	}

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

// Find finds a mood by the given ID
func (s *MoodService) Find(id uint) (*Mood, error) {
	mood, err := s.r.Find(id)
	if err != nil {
		return nil, err
	}
	return mood, nil
}

// DeleteForUser removes the stored mood by the given ID and user
func (s *MoodService) DeleteForUser(id uint, user *User) error {
	mood, err := s.FindForUser(id, user)
	if err != nil {
		return err
	}

	if err := s.r.Remove(id); err != nil {
		return err
	}

	// Add task to delete playlist in spotify
	if err := s.q.DeletePlaylist(user.ID, mood.PlaylistID); err != nil {
		return err
	}
	return nil
}

// CreatePlaylistForMood adds a new playlist in spotify for the given mood id
func (s *MoodService) CreatePlaylistForMood(name string, moodID uint, token *SpotifyToken) error {
	mood, err := s.Find(moodID)
	if err != nil {
		return err
	}

	id, err := s.spotify.CreatePlaylist(token, name)
	if err != nil {
		return err
	}

	mood.PlaylistID = id

	return s.r.Save(mood)
}
