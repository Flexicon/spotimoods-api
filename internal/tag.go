package internal

// Tag linking an artist to a mood
type Tag struct {
	ArtistID   string        `gorm:"primary_key;auto_increment:false" json:"artist_id"`
	MoodID     uint          `gorm:"primary_key;auto_increment:false" json:"mood_id"`
	ArtistData SpotifyArtist `gorm:"-" json:"artist_data,omitempty"`
}
