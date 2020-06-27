package model

// AddPlaylistPayload for queue messages
type AddPlaylistPayload struct {
	UserID uint   `json:"user_ID"`
	MoodID uint   `json:"mood_id"`
	Name   string `json:"name"`
}

// UpdatePlaylistPayload for queue messages
type UpdatePlaylistPayload struct {
	UserID     uint   `json:"user_ID"`
	PlaylistID string `json:"playlist_id"`
	Name       string `json:"name"`
}

// DeletePlaylistPayload for queue messages
type DeletePlaylistPayload struct {
	UserID     uint   `json:"user_ID"`
	PlaylistID string `json:"playlist_id"`
}
