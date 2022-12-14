package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type (
	QueueName    string
	ExchangeName string
	ExchangeKey  string
	Message      []byte
)

const (
	reconnectionAttemptsNum = 5
	reconnectionTimeout     = time.Second * 30
)

type Client struct {
	url      string
	conn     *amqp.Connection
	amqpChan *amqp.Channel

	connClients []ConnectionClient

	connCheckTimeout time.Duration
}

type Option func(c *Client) error

func Connect(url string, options ...Option) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to reconnect: %w", err)
	}

	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	c := &Client{
		url:      url,
		conn:     conn,
		amqpChan: amqpChan,

		connCheckTimeout: time.Second * 10, // default value
	}

	for _, opt := range options {
		if err = opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) Close() error {
	if c.amqpChan != nil {
		if err := c.amqpChan.Close(); err != nil {
			log.Errorf("Close amqp channel error: %v", err)
		}
	}

	if c.conn != nil && !c.conn.IsClosed() {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	return nil
}

func (c *Client) InitQueue(name QueueName) Queue {
	return &queue{
		name:   name,
		client: c,
	}
}

func (c *Client) InitExchange(name ExchangeName) Exchange {
	return &exchange{
		name:   name,
		client: c,
	}
}

func (c *Client) InitConsumer(queueName QueueName, options *ConsumerOptions, processor MessageProcessor) Consumer {
	return &consumer{
		client:           c,
		queue:            c.InitQueue(queueName),
		messageProcessor: processor,
		options:          options,
	}
}

func (c *Client) StartConsumers(ctx context.Context, consumers ...Consumer) error {
	for _, consumer := range consumers {
		if err := consumer.Start(ctx); err != nil {
			return fmt.Errorf("failed to start consumer: %w", err)
		}

		c.AddConnectionClient(consumer)
	}

	return nil
}

func (c *Client) AddConnectionClient(connClient ConnectionClient) {
	c.connClients = append(c.connClients, connClient)
}

func (c *Client) ListenConnectionAsync(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		if err := c.ListenConnection(ctx); err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()
}

func (c *Client) ListenConnection(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			if err := c.Close(); err != nil {
				return fmt.Errorf("close mq: %w", err)
			}

			return nil
		default:
			if err := c.checkConnection(ctx); err != nil {
				return fmt.Errorf("check mq connection: %w", err)
			}

			time.Sleep(time.Second * 10)
		}
	}
}

//nolint:goerr113
func (c *Client) checkConnection(ctx context.Context) error {
	if c.conn.IsClosed() {
		log.Warn("MQ connection lost")

		for i := 0; i < reconnectionAttemptsNum; i++ {
			time.Sleep(reconnectionTimeout)

			log.Info("Connecting to MQ... Attempt ", i+1)

			if err := c.reconnect(); err != nil {
				log.Errorf("Reconnect: %v", err)

				continue
			}

			for _, connClient := range c.connClients {
				if err := connClient.Reconnect(ctx); err != nil {
					log.Errorf("Reconnect for %+v: %v", connClient, err)

					continue
				}
			}

			log.Info("MQ connection established")

			return nil
		}

		return fmt.Errorf("failed to establish MQ connection")
	}

	return nil
}

func (c *Client) reconnect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("failed to reconnect: %w", err)
	}

	amqpChan, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	c.conn = conn
	c.amqpChan = amqpChan

	return nil
}

func publish(amqpChan *amqp.Channel, exchange ExchangeName, key ExchangeKey, body []byte) error {
	return publishWithConfig(amqpChan, exchange, key, body, PublishConfig{})
}

func publishWithConfig(amqpChan *amqp.Channel,
	exchange ExchangeName, key ExchangeKey, body []byte, cfg PublishConfig,
) error {
	headers := map[string]interface{}{}

	if cfg.MaxRetries != nil {
		headers[headerRemainingRetries] = *cfg.MaxRetries
	}

	err := amqpChan.Publish(string(exchange), string(key), false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
		Headers:      headers,
	})
	if err != nil {
		return fmt.Errorf("failed to publish a message to exchange: %w", err)
	}

	return nil
}

type ConnectionClient interface {
	Reconnect(ctx context.Context) error
}
