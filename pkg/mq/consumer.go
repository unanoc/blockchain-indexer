package mq

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/unanoc/blockchain-indexer/pkg/metrics"
	"github.com/unanoc/blockchain-indexer/pkg/nullable"
)

const headerRemainingRetries = "x-remaining-retries"

type consumer struct {
	client *Client

	queue            Queue
	messageProcessor MessageProcessor
	options          *ConsumerOptions

	messages <-chan amqp.Delivery
	stopChan chan struct{}
}

type Consumer interface {
	Start(ctx context.Context) error
	Reconnect(ctx context.Context) error
}

type MessageProcessor interface {
	Process(Message) error
}

// MessageProcessorFunc is an adapter to allow to use
// an ordinary functions as mq MessageProcessor.
type MessageProcessorFunc func(message Message) error

func (f MessageProcessorFunc) Process(m Message) error {
	return f(m)
}

func (c *consumer) Start(ctx context.Context) error {
	c.stopChan = make(chan struct{})

	var err error
	c.messages, err = c.messageChannel()
	if err != nil {
		return fmt.Errorf("failed to get message channel: %w", err)
	}
	for w := 1; w <= c.options.Workers; w++ {
		go c.consume(ctx)
	}

	log.Infof("Started %d MQ consumer workers for queue %s", c.options.Workers, c.queue.Name())

	return nil
}

func (c *consumer) Reconnect(ctx context.Context) error {
	c.messages = nil
	if c.stopChan != nil {
		close(c.stopChan)
	}

	err := c.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

//nolint:nestif,gocognit
func (c *consumer) consume(ctx context.Context) {
	queueName := string(c.queue.Name())

	for {
		select {
		case <-ctx.Done():
			log.Infof("Finished consuming queue %s", queueName)

			return
		case <-c.stopChan:
			log.Infof("Force stopped consuming queue %s", queueName)

			return
		case msg := <-c.messages:
			if msg.Body == nil {
				continue
			}

			err := c.process(msg.Body)
			if err != nil {
				log.Error(err)

				if c.options.RetryOnError {
					time.Sleep(c.options.RetryDelay)
					remainingRetries := c.getRemainingRetries(msg)

					switch {
					case remainingRetries > 0:
						if err := c.queue.PublishWithConfig(msg.Body, PublishConfig{
							MaxRetries: nullable.Int(int(remainingRetries - 1)),
						}); err != nil {
							log.Error(err)
						}
					case remainingRetries == 0:
						break
					default:
						if err := msg.Reject(true); err != nil {
							log.Error(err)
						}

						continue
					}
				}
			}

			if err := msg.Ack(false); err != nil {
				log.Error(err)
			}
		}
	}
}

func (c *consumer) process(body []byte) error {
	metric := c.options.PerformanceMetric
	if metric == nil {
		metric = &metrics.NullablePerformanceMetric{}
	}
	defer metric.Duration(metric.Start())

	if err := c.messageProcessor.Process(body); err != nil {
		metric.Failure()

		return fmt.Errorf("failed to process: %w", err)
	}

	metric.Success()

	return nil
}

func (c *consumer) messageChannel() (<-chan amqp.Delivery, error) {
	messageChannel, err := c.client.amqpChan.Consume(
		string(c.queue.Name()),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume a message: %w", err)
	}

	return messageChannel, nil
}

func (c *consumer) getRemainingRetries(delivery amqp.Delivery) int32 {
	remainingRetriesRaw, exists := delivery.Headers[headerRemainingRetries]
	if !exists {
		return int32(c.options.MaxRetries)
	}

	remainingRetries, ok := remainingRetriesRaw.(int32)
	if !ok {
		return int32(c.options.MaxRetries)
	}

	return remainingRetries
}
