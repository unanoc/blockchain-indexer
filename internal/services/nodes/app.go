package nodes

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/internal/services"
	"github.com/unanoc/blockchain-indexer/pkg/metrics"
	"github.com/unanoc/blockchain-indexer/pkg/service"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
)

type App struct {
	metricsPusher worker.Worker
	checker       worker.Worker
}

func NewApp() *App {
	services.InitConfig()
	services.InitLogging()
	services.InitSentry()
	services.InitDatabase()

	db, err := postgres.New(config.Default.Database.URL, config.Default.Database.Log)
	if err != nil {
		log.WithError(err).Fatal("Database init error")
	}

	prometheus := prometheus.NewPrometheus(config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem)
	prometheus.RegisterNodesMetrics()

	metricsPusher, err := metrics.InitDefaultMetricsPusher(
		config.Default.Prometheus.PushGateway.URL,
		config.Default.Prometheus.PushGateway.Key,
		fmt.Sprintf("%s_%s", config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem),
		config.Default.Prometheus.PushGateway.PushInterval,
	)
	if err != nil {
		log.WithError(err).Warn("Metrics pusher init error")
	}

	checker := NewWorker(db, prometheus)

	if config.Default.Nodes.InitNodes {
		if err := AddNodesListToDB(db); err != nil {
			log.WithError(err).Fatal("Nodes list adding error")
		}
	}

	return &App{
		metricsPusher: metricsPusher,
		checker:       checker,
	}
}

func (a *App) Run(ctx context.Context) {
	service.RunWithGracefulShutdown(ctx, func(ctx context.Context, wg *sync.WaitGroup) {
		if a.metricsPusher != nil {
			a.metricsPusher.Start(ctx, wg)
		}

		a.checker.Start(ctx, wg)
	})
}
