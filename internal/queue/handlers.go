package queue

import (
	"encoding/json"
	"log"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/flexicon/spotimoods-go/internal/queue/model"
	"github.com/streadway/amqp"
)

// Handler contains all handlers for every declared queue
type Handler struct {
	services *internal.ServiceProvider
}

// NewHandler constructor
func NewHandler(s *internal.ServiceProvider) *Handler {
	return &Handler{services: s}
}

func (h *Handler) handlePing(d amqp.Delivery) error {
	log.Printf("handling '%s': %s", pingQueue, d.Body)
	return nil
}

func (h *Handler) handleAddPlaylist(d amqp.Delivery) error {
	log.Printf("handling '%s': %s", addPlaylistQueue, d.Body)

	var payload model.AddPlaylistPayload
	if err := json.Unmarshal(d.Body, &payload); err != nil {
		return err
	}

	token, err := h.services.User().FindTokenForUser(payload.UserID)
	if err != nil {
		return err
	}

	if err := h.services.Mood().CreatePlaylistForMood(payload.Name, payload.MoodID, token); err != nil {
		return err
	}

	log.Printf(`Successfully created playlist for Mood ID %d, named: "%s"`, payload.MoodID, payload.Name)
	return nil
}

func (h *Handler) handleUpdatePlaylist(d amqp.Delivery) error {
	log.Printf("handling '%s': %s", updatePlaylistQueue, d.Body)

	var payload model.UpdatePlaylistPayload
	if err := json.Unmarshal(d.Body, &payload); err != nil {
		return err
	}

	token, err := h.services.User().FindTokenForUser(payload.UserID)
	if err != nil {
		return err
	}

	if err := h.services.Spotify().UpdatePlaylist(token, payload.PlaylistID, payload.Name); err != nil {
		return err
	}

	log.Printf("Successfully updated playlist: %s", payload.PlaylistID)

	return nil
}

func (h *Handler) handleDeletePlaylist(d amqp.Delivery) error {
	log.Printf("handling '%s': %s", deletePlaylistQueue, d.Body)

	var payload model.DeletePlaylistPayload
	if err := json.Unmarshal(d.Body, &payload); err != nil {
		return err
	}

	token, err := h.services.User().FindTokenForUser(payload.UserID)
	if err != nil {
		return err
	}

	if err := h.services.Spotify().DeletePlaylist(token, payload.PlaylistID); err != nil {
		return err
	}

	log.Printf("Successfully deleted playlist: %s", payload.PlaylistID)
	return nil
}
