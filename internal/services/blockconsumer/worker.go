package blockconsumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/pkg/mq"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
	"github.com/unanoc/blockchain-indexer/platform"
)

const (
	workerName = "block_consumer"
)

type Worker struct {
	log         *log.Entry
	kafka       *kafka.Reader
	prometheus  *prometheus.Prometheus
	API         platform.Platform
	txsExchange mq.Exchange
}

func NewWorker(txsExchange mq.Exchange, kafka *kafka.Reader,
	p *prometheus.Prometheus, pl platform.Platform,
) worker.Worker {
	w := &Worker{
		log: log.WithFields(log.Fields{
			"worker": workerName,
			"chain":  pl.Coin().Handle,
		}),
		txsExchange: txsExchange,
		kafka:       kafka,
		prometheus:  p,
		API:         pl,
	}

	opts := &worker.Options{
		Interval:        config.Default.BlockConsumer.Interval,
		RunImmediately:  true,
		RunConsequently: false,
	}

	return worker.NewWorkerBuilder(workerName, w.log, w.run).WithOptions(opts).Build()
}

func (w *Worker) run(ctx context.Context) error {
	chain := w.API.Coin().Handle

	w.prometheus.SetBlocksConsumerMetrics(w.kafka, chain)

	message, err := w.kafka.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch kafka message: %w", err)
	}

	txs, err := w.API.NormalizeRawBlock(message.Value)
	if err != nil {
		return fmt.Errorf("failed to normalize raw block: %w", err)
	}

	if err = w.publishToExchange(txs); err != nil {
		return fmt.Errorf("failed to publish to exchange: %w", err)
	}

	if err = w.kafka.CommitMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to commit kafka message (topic=%s offset=%d partition=%d): %w",
			message.Topic, message.Offset, message.Partition, err)
	}

	w.prometheus.SetBlocksConsumerTopicPartitionOffset(chain, message.Topic, message.Partition, message.Offset)

	log.WithFields(log.Fields{
		"chain":     w.API.Coin().Handle,
		"txs":       len(txs),
		"partition": message.Partition,
		"offset":    message.Offset,
	}).Info("Transactions have been consumed")

	return nil
}

func (w *Worker) publishToExchange(txs types.Txs) error {
	if len(txs) == 0 {
		return nil
	}

	logFields := log.Fields{"chain": w.API.Coin().Handle, "txs": txs[:1]}

	body, err := json.Marshal(txs)
	if err != nil {
		log.WithFields(logFields).Error(err)

		return fmt.Errorf("failed to marshal json: %w", err)
	}

	if err = w.txsExchange.Publish(body); err != nil {
		log.WithFields(logFields).Error(err)

		return fmt.Errorf("failed to publish txs to rabbit exchange: %w", err)
	}

	return nil
}
