package mq

import "fmt"

type exchange struct {
	name   ExchangeName
	client *Client
}

type Exchange interface {
	Declare(kind string) error
	Bind(queues []Queue) error
	BindWithKey(queues []Queue, key ExchangeKey) error
	Publish(body []byte) error
	PublishWithKey(body []byte, key ExchangeKey) error
}

func (e *exchange) Declare(kind string) error {
	err := e.client.amqpChan.ExchangeDeclare(string(e.name), kind, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	return nil
}

func (e *exchange) Bind(queues []Queue) error {
	for _, q := range queues {
		err := e.client.amqpChan.QueueBind(string(q.Name()), "", string(e.name), false, nil)
		if err != nil {
			return fmt.Errorf("failed to bind a queue: %w", err)
		}
	}

	return nil
}

func (e *exchange) BindWithKey(queues []Queue, key ExchangeKey) error {
	for _, q := range queues {
		err := e.client.amqpChan.QueueBind(string(q.Name()), string(key), string(e.name), false, nil)
		if err != nil {
			return fmt.Errorf("failed to bind a queue: %w", err)
		}
	}

	return nil
}

func (e *exchange) Publish(body []byte) error {
	return publish(e.client.amqpChan, e.name, "", body)
}

func (e *exchange) PublishWithKey(body []byte, key ExchangeKey) error {
	return publish(e.client.amqpChan, e.name, key, body)
}
