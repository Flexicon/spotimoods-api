package queue

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

const (
	pingQueue           = "ping"
	addPlaylistQueue    = "add_playlist"
	updatePlaylistQueue = "update_playlist"
	deletePlaylistQueue = "delete_playlist"
)

var durableQueues = []string{addPlaylistQueue, updatePlaylistQueue, deletePlaylistQueue}

// Service to manage working with the queue
type Service struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// Setup declares the necessary queue topology and returns a new queue service
func Setup() (*Service, error) {
	conn, err := amqp.Dial(viper.GetString("rabbit.uri"))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Failed to open a channel: %v", err)
	}

	if _, err = ch.QueueDeclare(
		pingQueue, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	); err != nil {
		return nil, fmt.Errorf("Failed to declare ping queue: %v", err)
	}

	for _, name := range durableQueues {
		if _, err = ch.QueueDeclare(
			name,  // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		); err != nil {
			return nil, fmt.Errorf("Failed to declare %s queue: %v", name, err)
		}
	}

	return &Service{
		ch:   ch,
		conn: conn,
	}, nil
}

// Listen sets up consumers and begins listening for messages
func Listen(s *Service, h *Handler) error {
	defer s.conn.Close()
	defer s.ch.Close()

	pings, err := s.ch.Consume(
		pingQueue, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("Failed to register consumer: %v", err)
	}

	addPlaylistMsgs, err := s.ch.Consume(
		addPlaylistQueue, // queue
		"",               // consumer
		false,            // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		return fmt.Errorf("Failed to register consumer: %v", err)
	}

	updatePlaylistMsgs, err := s.ch.Consume(
		updatePlaylistQueue, // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		return fmt.Errorf("Failed to register consumer: %v", err)
	}

	deletePlaylistMsgs, err := s.ch.Consume(
		deletePlaylistQueue, // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		return fmt.Errorf("Failed to register consumer: %v", err)
	}

	untilErr := make(chan error)

	go handleMessages(pings, h.handlePing)
	go handleMessages(addPlaylistMsgs, h.handleAddPlaylist)
	go handleMessages(updatePlaylistMsgs, h.handleUpdatePlaylist)
	go handleMessages(deletePlaylistMsgs, h.handleDeletePlaylist)

	return <-untilErr
}

func handleMessages(msgs <-chan amqp.Delivery, handler func(d amqp.Delivery) error) {
	for d := range msgs {
		if err := handler(d); err != nil {
			log.Printf("message rejected: %v", err)
			if err := d.Reject(false); err != nil {
				log.Printf("failed to reject message: %v", err)
			}
		}
		if err := d.Ack(false); err != nil {
			log.Printf("failed to acknowledge message: %v", err)
		}
	}
}
