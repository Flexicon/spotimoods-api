package queue

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

func (s *Service) publishJSON(queue string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := s.ch.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	); err != nil {
		return err
	}

	return nil
}
