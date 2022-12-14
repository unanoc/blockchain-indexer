package blockconsumer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/rabbit"
	"github.com/unanoc/blockchain-indexer/internal/services"
	"github.com/unanoc/blockchain-indexer/pkg/metrics"
	"github.com/unanoc/blockchain-indexer/pkg/mq"
	"github.com/unanoc/blockchain-indexer/pkg/service"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
	"github.com/unanoc/blockchain-indexer/platform"
)

type App struct {
	metricsPusher worker.Worker
	workers       []worker.Worker
}

func NewApp() *App {
	services.InitConfig()
	services.InitLogging()
	services.InitSentry()
	services.InitRabbitMQ()

	platforms := platform.InitPlatforms()

	rabbitmq, err := mq.Connect(config.Default.RabbitMQ.URL)
	if err != nil {
		log.WithError(err).Fatal("RabbitMQ init error")
	}

	prometheus := prometheus.NewPrometheus(config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem)
	prometheus.RegisterBlocksConsumerMetrics()

	metricsPusher, err := metrics.InitDefaultMetricsPusher(
		config.Default.Prometheus.PushGateway.URL,
		config.Default.Prometheus.PushGateway.Key,
		fmt.Sprintf("%s_%s", config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem),
		config.Default.Prometheus.PushGateway.PushInterval,
	)
	if err != nil {
		log.WithError(err).Warn("Metrics pusher init error")
	}

	txsExchange := rabbitmq.InitExchange(rabbit.ExchangeTransactionsParsed)

	workers := make([]worker.Worker, 0, len(platforms))
	for _, pl := range platforms {
		kafka := kafka.NewReader(kafka.ReaderConfig{
			Brokers:       strings.Split(config.Default.Kafka.Brokers, ","),
			MaxAttempts:   config.Default.Kafka.MaxAttempts,
			Topic:         fmt.Sprintf("%s%s", config.Default.Kafka.BlocksTopicPrefix, pl.Coin().Handle),
			GroupID:       pl.Coin().Handle,
			StartOffset:   kafka.FirstOffset,
			RetentionTime: config.Default.Kafka.RetentionTime,
		})

		workers = append(workers, NewWorker(txsExchange, kafka, prometheus, pl))
	}

	return &App{
		metricsPusher: metricsPusher,
		workers:       workers,
	}
}

func (a *App) Run(ctx context.Context) {
	service.RunWithGracefulShutdown(ctx, func(ctx context.Context, wg *sync.WaitGroup) {
		if a.metricsPusher != nil {
			a.metricsPusher.Start(ctx, wg)
		}

		for _, worker := range a.workers {
			go worker.Start(ctx, wg)
		}
	})
}
