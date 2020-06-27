package internal

// QueueService manages the message queue
type QueueService interface {
	// AddPlaylist publishes a new message to the add_playlist queue
	AddPlaylist(mood *Mood) error
	// UpdatePlaylist publishes a new message to the update_playlist queue
	UpdatePlaylist(mood *Mood) error
	// DeletePlaylist publishes a new message to the delete_playlist queue
	DeletePlaylist(userID uint, playlistID string) error
}
