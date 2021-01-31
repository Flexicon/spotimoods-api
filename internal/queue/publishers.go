package queue

import (
	"github.com/flexicon/spotimoods-go/internal"
	"github.com/flexicon/spotimoods-go/internal/queue/model"
)

// Ping publishes a new message to the ping queue
func (s *Service) Ping(msg string) error {
	payload := model.PingPayload{Msg: msg}

	return s.publishJSON(pingQueue, payload)
}

// AddPlaylist publishes a new message to the add_playlist queue
func (s *Service) AddPlaylist(mood *internal.Mood) error {
	payload := model.AddPlaylistPayload{
		UserID: mood.UserID,
		MoodID: mood.ID,
		Name:   mood.Name,
	}

	return s.publishJSON(addPlaylistQueue, payload)
}

// UpdatePlaylist publishes a new message to the update_playlist queue
func (s *Service) UpdatePlaylist(mood *internal.Mood) error {
	payload := model.UpdatePlaylistPayload{
		UserID:     mood.UserID,
		PlaylistID: mood.PlaylistID,
		Name:       mood.Name,
	}

	return s.publishJSON(updatePlaylistQueue, payload)
}

// DeletePlaylist publishes a new message to the delete_playlist queue
func (s *Service) DeletePlaylist(userID uint, playlistID string) error {
	payload := model.DeletePlaylistPayload{UserID: userID, PlaylistID: playlistID}

	return s.publishJSON(deletePlaylistQueue, payload)
}
