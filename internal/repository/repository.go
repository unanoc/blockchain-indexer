package repository

import (
	"context"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
)

type Storage interface {
	// Block trackers
	GetBlockTracker(ctx context.Context, chain string) (*models.BlockTracker, error)
	UpsertBlockTracker(ctx context.Context, chain string, height int64) error

	// Transactions
	InsertTransactions(ctx context.Context, txs []models.Transaction) error
	GetTransactions(ctx context.Context, chain string, page, limit int, recent bool) ([]models.Transaction, error)
	GetTransactionsByAddress(ctx context.Context,
		chain, address string, page, limit int, recent bool) ([]models.Transaction, error)
	GetTransactionByHash(ctx context.Context, chain, hash string) (*models.Transaction, error)
	GetTransactionTotalCount(ctx context.Context, chain string) (int64, error)
	GetTransactionByAddressTotalCount(ctx context.Context, chain, address string) (int64, error)

	// Nodes
	InsertNodes(ctx context.Context, nodes []models.Node) error
	GetNodes(ctx context.Context) ([]models.Node, error)
	GetNodesByChain(ctx context.Context, chain string) ([]models.Node, error)
}
