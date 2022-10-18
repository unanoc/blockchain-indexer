package transactions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/internal/services/api/httperr"
)

var ErrTxDoesNotExist = errors.New("transaction does not exist")

type Controller struct {
	db repository.Storage
}

func NewController(db repository.Storage) *Controller {
	return &Controller{db: db}
}

func (i *Controller) GetTransactions(ctx context.Context, chain string,
	page, limit int, recent bool,
) (*TxsResp, *httperr.Error) {
	txs, err := i.db.GetTransactions(ctx, chain, page, limit, recent)
	if err != nil {
		log.WithError(err).Error("Getting txs error")

		return nil, httperr.ErrInternalServer
	}

	transactions, err := ToTxs(txs)
	if err != nil {
		log.WithError(err).Error("Txs normalizing error")

		return nil, httperr.ErrInternalServer
	}

	totalCount, err := i.db.GetTransactionTotalCount(ctx, chain)
	if err != nil {
		log.WithError(err).Error("Getting of txs count error")

		return nil, httperr.ErrInternalServer
	}

	return &TxsResp{
		StatusCode: http.StatusOK,
		TotalCount: totalCount,
		TotalPages: int(math.Ceil(float64(totalCount) / float64(limit))),
		PageNumber: page,
		Limit:      limit,
		Txs:        transactions,
	}, nil
}

func (i *Controller) GetTransactionByHash(ctx context.Context, chain, hash string) (*TxResp, *httperr.Error) {
	tx, err := i.db.GetTransactionByHash(ctx, chain, hash)
	if err != nil {
		if postgres.IsErrNotFound(err) {
			return nil, httperr.NewError(http.StatusNotFound, ErrTxDoesNotExist.Error())
		}

		log.WithError(err).Error("Getting tx error")

		return nil, httperr.ErrInternalServer
	}

	var metadata interface{}
	if err := json.Unmarshal(tx.Metadata.RawMessage, &metadata); err != nil {
		log.WithError(err).Error("Tx normalization error")

		return nil, httperr.ErrInternalServer
	}

	return &TxResp{
		StatusCode: http.StatusOK,
		Tx: Tx{
			Hash:      tx.Hash,
			Chain:     tx.Chain,
			Height:    tx.Block,
			From:      tx.Sender,
			To:        tx.Recipient,
			Status:    string(tx.Status),
			Type:      string(tx.Type),
			Sequence:  tx.Sequence,
			Fee:       tx.Fee,
			Data:      metadata,
			Timestamp: tx.BlockCreatedAt,
		},
	}, nil
}

func (i *Controller) GetTransactionsByUser(ctx context.Context, chain string,
	address string, page, limit int, recent bool,
) (*TxsResp, *httperr.Error) {
	txs, err := i.db.GetTransactionsByAddress(ctx, chain, address, page, limit, recent)
	if err != nil {
		log.WithError(err).Error("Getting txs by address error")

		return nil, httperr.ErrInternalServer
	}

	transactions, err := ToTxs(txs)
	if err != nil {
		log.WithError(err).Error("Txs normalizing error")

		return nil, httperr.ErrInternalServer
	}

	totalCount, err := i.db.GetTransactionByAddressTotalCount(ctx, chain, address)
	if err != nil {
		log.WithError(err).Error("Getting of user txs count error")

		return nil, httperr.ErrInternalServer
	}

	return &TxsResp{
		StatusCode: http.StatusOK,
		TotalCount: totalCount,
		TotalPages: int(math.Ceil(float64(totalCount) / float64(limit))),
		PageNumber: page,
		Limit:      limit,
		Txs:        transactions,
	}, nil
}

func ToTxs(txs []models.Transaction) ([]Tx, error) {
	transactions := make([]Tx, 0, len(txs))

	for _, tx := range txs {
		var metadata interface{}
		if err := json.Unmarshal(tx.Metadata.RawMessage, &metadata); err != nil {
			return nil, fmt.Errorf("failed to normalize db.Transaction to Tx: %w", err)
		}

		transactions = append(transactions, Tx{
			Hash:      tx.Hash,
			Chain:     tx.Chain,
			Height:    tx.Block,
			From:      tx.Sender,
			To:        tx.Recipient,
			Status:    string(tx.Status),
			Type:      string(tx.Type),
			Sequence:  tx.Sequence,
			Fee:       tx.Fee,
			Data:      metadata,
			Timestamp: tx.BlockCreatedAt,
		})
	}

	return transactions, nil
}
