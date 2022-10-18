package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	httplib "github.com/unanoc/blockchain-indexer/pkg/http"
)

// API contains generic methods that should be used by other API interfaces.
type API interface {
	Setup(router *gin.Engine)
}

func NewRouter(db repository.Storage, p *prometheus.Prometheus) http.Handler {
	var router *gin.Engine

	if config.Default.Gin.Mode == gin.DebugMode {
		router = gin.Default()
	} else {
		router = gin.New()
	}

	// Setup service routers
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup middlewares
	router.Use(httplib.CORSMiddleware())
	router.Use(prometheus.GinMetricsMiddleware(p))

	// Setup API routes
	NewTransactionsAPI(db).Setup(router)

	return router
}

func (api *TransactionsAPI) Setup(router *gin.Engine) {
	router.GET("/api/v1/transactions", api.GetTransactions)
	router.GET("/api/v1/transactions/:hash", api.GetTransactionByHash)
	router.GET("/api/v1/transactions/user", api.GetTransactionsByUser)
}
