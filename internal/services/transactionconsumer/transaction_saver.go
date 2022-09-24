package transactionconsumer

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/pkg/mq"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

const serviceName = "transaction_saver"

type TransactionSaver struct {
	db         *postgres.Database
	prometheus *prometheus.Prometheus
}

func NewTransactionSaver(db *postgres.Database, p *prometheus.Prometheus) mq.MessageProcessor {
	return &TransactionSaver{
		db:         db,
		prometheus: p,
	}
}

func (ts *TransactionSaver) Process(message mq.Message) error {
	var txs types.Txs

	if err := json.Unmarshal(message, &txs); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	log.WithFields(log.Fields{
		"service": serviceName,
		"txs":     len(txs),
	}).Info("Consumed")

	if len(txs) == 0 {
		return nil
	}

	chain := txs[0].Chain

	normalizedTxs, err := models.NormalizeTransactions(txs, chain)
	if err != nil {
		return fmt.Errorf("failed to normalized txs: %w", err)
	}

	if err = ts.db.InsertTransactions(context.Background(), normalizedTxs); err != nil {
		return fmt.Errorf("failed to insert txs: %w", err)
	}

	ts.prometheus.SetParsedTxs(chain, len(normalizedTxs))

	return nil
}
