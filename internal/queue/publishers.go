package queue

import (
	"github.com/flexicon/spotimoods-go/internal"
	"github.com/flexicon/spotimoods-go/internal/queue/model"
)

// Ping publishes a new message to the ping queue
func (s *Service) Ping(msg string) error {
	payload := model.PingPayload{Msg: msg}

	if err := s.publishJSON(pingQueue, payload); err != nil {
		return err
	}
	return nil
}

// AddPlaylist publishes a new message to the add_playlist queue
func (s *Service) AddPlaylist(mood *internal.Mood) error {
	payload := model.AddPlaylistPayload{
		UserID: mood.UserID,
		MoodID: mood.ID,
		Name:   mood.Name,
	}

	if err := s.publishJSON(addPlaylistQueue, payload); err != nil {
		return err
	}
	return nil
}

// UpdatePlaylist publishes a new message to the update_playlist queue
func (s *Service) UpdatePlaylist(mood *internal.Mood) error {
	payload := model.UpdatePlaylistPayload{
		UserID:     mood.UserID,
		PlaylistID: mood.PlaylistID,
		Name:       mood.Name,
	}

	if err := s.publishJSON(updatePlaylistQueue, payload); err != nil {
		return err
	}
	return nil
}

// DeletePlaylist publishes a new message to the delete_playlist queue
func (s *Service) DeletePlaylist(userID uint, playlistID string) error {
	payload := model.DeletePlaylistPayload{UserID: userID, PlaylistID: playlistID}

	if err := s.publishJSON(deletePlaylistQueue, payload); err != nil {
		return err
	}
	return nil
}
