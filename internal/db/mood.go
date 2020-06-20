package db

import (
	"github.com/flexicon/spotimoods-go/internal"
	"github.com/jinzhu/gorm"
)

// MoodRepository for interacting with mood data in the DB
type MoodRepository struct {
	db *gorm.DB
}

// Find mood by ID
func (r *MoodRepository) Find(id uint) (*internal.Mood, error) {
	var mood internal.Mood
	query := r.db.First(&mood, id)
	if query.RecordNotFound() {
		return nil, internal.ErrNotFound
	}

	return &mood, query.Error
}

// FindByIDAndUser if it exists
func (r *MoodRepository) FindByIDAndUser(id uint, user *internal.User) (*internal.Mood, error) {
	var mood internal.Mood
	query := r.db.Where("id = ? AND user_id = ?", id, user.ID).First(&mood)
	if query.RecordNotFound() {
		return nil, internal.ErrNotFound
	}

	return &mood, query.Error
}

// Remove mood by ID
func (r *MoodRepository) Remove(id uint) error {
	query := r.db.Delete(internal.Mood{ID: id})
	if query.RecordNotFound() {
		return internal.ErrNotFound
	}

	return query.Error
}

// Save upserts the given user into the DB
func (r *MoodRepository) Save(mood *internal.Mood) error {
	if r.db.NewRecord(mood) {
		return r.db.Create(&mood).Error
	}
	return r.db.Save(&mood).Error
}

// FindByUser all moods for a given user
func (r *MoodRepository) FindByUser(user *internal.User) ([]*internal.Mood, error) {
	var moods []*internal.Mood
	err := r.db.Where("user_id = ?", user.ID).Find(&moods).Error

	return moods, err
}
