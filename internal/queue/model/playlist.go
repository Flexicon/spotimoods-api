package model

// AddPlaylistPayload for queue messages
type AddPlaylistPayload struct {
	MoodID uint `json:"mood_id"`
}

// UpdatePlaylistPayload for queue messages
type UpdatePlaylistPayload struct {
	MoodID uint   `json:"mood_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

// DeletePlaylistPayload for queue messages
type DeletePlaylistPayload struct {
	UserID     uint   `json:"user_ID"`
	PlaylistID string `json:"playlist_id"`
}
