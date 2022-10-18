package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
)

const batchSize = 1000

func (d *Database) InsertTransactions(ctx context.Context, txs []models.Transaction) error {
	if err := d.Gorm.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&txs, batchSize).Error; err != nil {
		return fmt.Errorf("failed to insert txs: %w", err)
	}

	return nil
}

func (d *Database) GetTransactions(ctx context.Context, chain string, page, limit int, recent bool,
) ([]models.Transaction, error) {
	txs := make([]models.Transaction, 0)

	orderBy := "asc"
	if recent {
		orderBy = "desc"
	}

	if err := d.Gorm.
		WithContext(ctx).
		Where("chain = ?", chain).
		Order("block_created_at " + orderBy).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&txs).Error; err != nil {
		return nil, fmt.Errorf("failed to find txs by address: %w", err)
	}

	return txs, nil
}

func (d *Database) GetTransactionsByAddress(ctx context.Context,
	chain, address string, page, limit int, recent bool,
) ([]models.Transaction, error) {
	txs := make([]models.Transaction, 0)

	orderBy := "asc"
	if recent {
		orderBy = "desc"
	}

	if err := d.Gorm.
		WithContext(ctx).
		Where("chain = ?", chain).
		Where("sender = ? or recipient = ?", address, address).
		Order("block_created_at " + orderBy).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&txs).Error; err != nil {
		return nil, fmt.Errorf("failed to find txs by address: %w", err)
	}

	return txs, nil
}

func (d *Database) GetTransactionByHash(ctx context.Context, chain, hash string) (*models.Transaction, error) {
	var tx models.Transaction

	if err := d.Gorm.
		WithContext(ctx).
		Where("chain = ?", chain).
		Where("hash = ?", hash).
		First(&tx).Error; err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	return &tx, nil
}

func (d *Database) GetTransactionTotalCount(ctx context.Context, chain string) (int64, error) {
	var totalCount int64

	if err := d.Gorm.
		WithContext(ctx).
		Model(&models.Transaction{}).
		Where("chain = ?", chain).
		Count(&totalCount).Error; err != nil {
		return totalCount, fmt.Errorf("failed to get count of txs: %w", err)
	}

	return totalCount, nil
}

func (d *Database) GetTransactionByAddressTotalCount(ctx context.Context, chain, address string) (int64, error) {
	var totalCount int64

	if err := d.Gorm.
		WithContext(ctx).
		Model(&models.Transaction{}).
		Where("chain = ?", chain).
		Where("sender = ? or recipient = ?", address, address).
		Count(&totalCount).Error; err != nil {
		return totalCount, fmt.Errorf("failed to get count of txs for address: %w", err)
	}

	return totalCount, nil
}
