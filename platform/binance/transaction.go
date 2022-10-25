package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

var (
	ErrBlocktimeOverflow = errors.New("BlockTime causes number overflow")
	ErrFromAddrNotString = errors.New("FromAddr is not a string")
	ErrToAddrNotString   = errors.New("ToAddr is not a string")
	ErrAmountNil         = errors.New("transfer amount cannot be nil")
)

func (p *Platform) NormalizeRawBlock(rawBlock []byte) (types.Txs, error) {
	var block Block
	if err := json.Unmarshal(rawBlock, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return p.NormalizeBlock(block), nil
}

func (p *Platform) NormalizeBlock(block Block) types.Txs {
	normalizedTxs := make(types.Txs, 0, len(block.Txs))

	for _, t := range block.Txs {
		logger := log.WithField("tx_id", t.TxHash).WithField("coin", p.Coin())

		token := t.TxAsset
		if t.TxAsset == coin.Binance().Symbol {
			token = ""
		}

		assetID := coin.Binance().TokenAssetID(token)

		switch t.TxType {
		case Transfer:
			if len(t.Data) == 0 {
				normalizedTx, err := p.normalizeTransferTransaction(t, assetID)
				if err != nil {
					logger.WithError(err).Error("Error while normalizing tranfer tx")

					continue
				}

				normalizedTxs = append(normalizedTxs, normalizedTx)
			}
		case Delegate, Undelegate:
			normalizedTx, err := p.normalizeDelegationTransaction(t, assetID)
			if err != nil {
				logger.WithError(err).Errorf("Error while normalizing delegation tx")

				continue
			}

			normalizedTxs = append(normalizedTxs, normalizedTx)
		}
	}

	return normalizedTxs
}

func (p *Platform) normalizeBaseOfTx(t Tx, assetID coin.AssetID) (types.Tx, error) {
	tx := types.Tx{
		Hash:     t.TxHash,
		Chain:    p.Coin().Handle,
		Fee:      normalizeFee(t.TxFee, assetID),
		Block:    uint64(t.BlockHeight),
		Status:   types.StatusSuccess,
		Sequence: uint64(t.Sequence),
	}

	blockTime := t.BlockTime / 1000
	if blockTime > math.MaxInt64 {
		return tx, ErrBlocktimeOverflow
	}
	tx.BlockCreatedAt = int64(blockTime)

	return tx, nil
}

func (p *Platform) normalizeTransferTransaction(t Tx, assetID coin.AssetID) (types.Tx, error) {
	tx, err := p.normalizeBaseOfTx(t, assetID)
	if err != nil {
		return tx, err
	}

	tx.Type = types.TxTransfer
	tx.Memo = t.Memo

	var ok bool
	if tx.From, ok = t.FromAddr.(string); !ok {
		return tx, ErrFromAddrNotString
	}
	if tx.To, ok = t.ToAddr.(string); !ok {
		return tx, ErrToAddrNotString
	}
	if t.Amount == nil {
		return tx, ErrAmountNil
	}

	tx.Metadata = &types.Transfer{
		Asset:  assetID,
		Amount: types.Amount(strconv.FormatUint(*t.Amount, 10)),
	}

	return tx, nil
}

func (p *Platform) normalizeDelegationTransaction(t Tx, assetID coin.AssetID) (types.Tx, error) {
	tx, err := p.normalizeBaseOfTx(t, assetID)
	if err != nil {
		return tx, err
	}

	var ok bool
	if tx.From, ok = t.FromAddr.(string); !ok {
		return tx, ErrFromAddrNotString
	}

	var data DelegationData
	if err := json.Unmarshal([]byte(t.Data), &data); err != nil {
		return types.Tx{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	tx.To = data.ValidatorAddress

	var amount uint64
	if t.TxType == Undelegate {
		tx.Type = types.TxStakeUndelegate
		tx.Memo = t.Memo
		amount = data.Amount.Amount
	} else {
		tx.Type = types.TxStakeDelegate
		tx.Memo = t.Memo
		amount = data.Delegation.Amount
	}

	tx.Metadata = &types.Transfer{
		Asset:  assetID,
		Amount: types.Amount(strconv.FormatUint(amount, 10)),
	}

	return tx, nil
}

func normalizeFee(amount uint64, asset coin.AssetID) types.Fee {
	return types.Fee{
		Asset:  asset,
		Amount: types.Amount(strconv.FormatUint(amount, 10)),
	}
}
