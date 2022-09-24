package transactionconsumer

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/rabbit"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/internal/services"
	"github.com/unanoc/blockchain-indexer/pkg/metrics"
	"github.com/unanoc/blockchain-indexer/pkg/mq"
	"github.com/unanoc/blockchain-indexer/pkg/service"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
)

type App struct {
	metricsPusher worker.Worker
	rabbitmq      *mq.Client
	consumers     []mq.Consumer
}

func NewApp() *App {
	services.InitConfig()
	services.InitLogging()
	services.InitSentry()
	services.InitDatabase()
	services.InitRabbitMQ()

	db, err := postgres.New(config.Default.Database.URL, config.Default.Database.Log)
	if err != nil {
		log.WithError(err).Fatal("Database init error")
	}

	rabbitmq, err := mq.Connect(config.Default.RabbitMQ.URL)
	if err != nil {
		log.WithError(err).Fatal("RabbitMQ init error")
	}

	prometheus := prometheus.NewPrometheus(config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem)
	prometheus.RegisterTrasactionConsumerMetrics()

	metricsPusher, err2 := metrics.InitDefaultMetricsPusher(
		config.Default.Prometheus.PushGateway.URL,
		config.Default.Prometheus.PushGateway.Key,
		fmt.Sprintf("%s_%s", config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem),
		config.Default.Prometheus.PushGateway.PushInterval,
	)
	if err2 != nil {
		log.WithError(err2).Warn("Metrics pusher init error")
	}

	consumers := []mq.Consumer{
		initTransactionSaver(db, rabbitmq, prometheus),
	}

	return &App{
		metricsPusher: metricsPusher,
		rabbitmq:      rabbitmq,
		consumers:     consumers,
	}
}

func (a *App) Run(ctx context.Context) {
	service.RunWithGracefulShutdown(ctx, func(ctx context.Context, wg *sync.WaitGroup) {
		if a.metricsPusher != nil {
			a.metricsPusher.Start(ctx, wg)
		}

		if err := a.rabbitmq.StartConsumers(ctx, a.consumers...); err != nil {
			log.WithError(err).Fatal("Rabbit MQ consumers starting error")
		}

		a.rabbitmq.ListenConnectionAsync(ctx, wg)
	})
}

func initTransactionSaver(db *postgres.Database, rmq *mq.Client, p *prometheus.Prometheus) mq.Consumer {
	return rmq.InitConsumer(rabbit.QueueTransactionsSave,
		rabbit.NewConsumerOptions(config.Default.TransactionConsumer.Workers), NewTransactionSaver(db, p))
}
