package api

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/internal/services"
	"github.com/unanoc/blockchain-indexer/internal/services/api/handlers"
	"github.com/unanoc/blockchain-indexer/pkg/http"
	"github.com/unanoc/blockchain-indexer/pkg/service"
)

type App struct {
	server http.Server
}

func NewApp() *App {
	services.InitConfig()
	services.InitLogging()
	services.InitSentry()
	services.InitDatabase()
	services.InitSwaggerInfo()

	db, err := postgres.New(config.Default.Database.URL, config.Default.Database.Log)
	if err != nil {
		log.WithError(err).Fatal("Database init error")
	}

	prometheus := prometheus.NewPrometheus(config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem)
	prometheus.RegisterAPIMetrics()

	router := handlers.NewRouter(db, prometheus)
	server := http.NewHTTPServer(router, config.Default.Port)

	return &App{server: server}
}

func (a *App) Run(ctx context.Context) {
	service.RunWithGracefulShutdown(ctx, a.server.Run)
}
