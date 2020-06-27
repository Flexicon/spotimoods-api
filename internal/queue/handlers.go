package queue

import (
	"log"

	"github.com/flexicon/spotimoods-go/internal"
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

	return nil
}

func (h *Handler) handleUpdatePlaylist(d amqp.Delivery) error {
	log.Printf("handling '%s': %s", updatePlaylistQueue, d.Body)

	return nil
}

func (h *Handler) handleDeletePlaylist(d amqp.Delivery) error {
	log.Printf("handling '%s': %s", deletePlaylistQueue, d.Body)

	return nil
}
