package cosmos

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/pkg/primitives/numbers"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

var ErrUnsupportedType = errors.New("unsupported message value type")

func (p *Platform) NormalizeRawBlock(rawBlock []byte) (types.Txs, error) {
	var block TxPage
	if err := json.Unmarshal(rawBlock, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return p.NormalizeBlock(block), nil
}

func (p *Platform) NormalizeBlock(block TxPage) types.Txs {
	srcTxs := block.TxResponses
	txMap := make(map[string]bool)
	txs := make(types.Txs, 0, len(srcTxs))

	for _, srcTx := range srcTxs {
		if _, seen := txMap[srcTx.TxHash]; seen {
			continue
		}

		normalizedTx, err := p.normalizeTx(srcTx)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"tx_id": srcTx.TxHash,
				"chain": p.Coin().Handle,
			}).Debug("Cannot normalize transaction")

			continue
		}

		txMap[srcTx.TxHash] = true

		if normalizedTx.Type != "" {
			txs = append(txs, normalizedTx)
		}
	}

	return txs
}

//nolint:nestif
func (p *Platform) normalizeTx(srcTx TxResponse) (tx types.Tx, err error) {
	date, err := time.Parse(time.RFC3339, srcTx.Date)
	if err != nil {
		return types.Tx{}, fmt.Errorf("failed to parse tx date: %w", err)
	}

	block, err := strconv.ParseUint(srcTx.Height, 10, 64)
	if err != nil {
		return types.Tx{}, fmt.Errorf("failed to parse block num: %w", err)
	}

	status := types.StatusSuccess
	if srcTx.Code > 0 {
		status = types.StatusError
	}

	fee := "0"
	feeToken := ""
	if fees := srcTx.Tx.AuthInfo.Fee.Amount; len(fees) > 0 {
		if qty := fees[0].Amount; len(qty) > 0 && qty != fee {
			fee, err = numbers.DecimalToSatoshis(qty)
			if err != nil {
				return types.Tx{}, fmt.Errorf("cannot convert decimal to satoshis: %w", err)
			}

			feeToken = fees[0].Denom

			if p.isCoin(feeToken) {
				feeToken = ""
			}
		}
	}

	var sequence uint64
	if len(srcTx.Tx.AuthInfo.SignerInfos) >= 1 {
		sequenceStr := srcTx.Tx.AuthInfo.SignerInfos[0].Sequence

		sequence, err = strconv.ParseUint(sequenceStr, 10, 64)
		if err != nil {
			return types.Tx{}, fmt.Errorf("failed to parse sequence: %w", err)
		}
	}

	tx = types.Tx{
		Hash:           srcTx.TxHash,
		Chain:          p.Coin().Handle,
		BlockCreatedAt: date.Unix(),
		Status:         status,
		Block:          block,
		Sequence:       sequence,
		Fee: types.Fee{
			Asset:  p.Coin().TokenAssetID(feeToken),
			Amount: types.Amount(fee),
		},
	}

	if len(srcTx.Tx.Body.Messages) == 0 {
		return tx, nil
	}

	msg := srcTx.Tx.Body.Messages[0]
	switch value := msg.MessageValue.(type) {
	case MessageValueSend:
		p.fillSend(&tx, srcTx, value)

		return tx, nil
	case MessageValueDelegate:
		p.fillDelegate(&tx, srcTx, value)

		return tx, nil
	default:
		return tx, ErrUnsupportedType
	}
}

func (p *Platform) fillSend(tx *types.Tx, srcTx TxResponse, msg MessageValueSend) {
	if len(msg.Amount) == 0 {
		return
	}

	amount, err := numbers.DecimalToSatoshis(msg.Amount[0].Amount)
	if err != nil {
		return
	}

	tx.From = msg.FromAddr
	tx.To = msg.ToAddr

	token := p.normalizeToken(msg.Amount[0].Denom)
	assetID := p.Coin().TokenAssetID(token)

	tx.Type = types.TxTransfer
	tx.Memo = srcTx.Tx.Body.Memo
	tx.Metadata = &types.Transfer{
		Amount: types.Amount(amount),
		Asset:  assetID,
	}
}

//nolint:exhaustive
func (p *Platform) fillDelegate(tx *types.Tx, srcTx TxResponse, msg MessageValueDelegate) {
	amount := ""
	if len(msg.Amount.Amount) > 0 {
		var err error
		amount, err = numbers.DecimalToSatoshis(msg.Amount.Amount)
		if err != nil {
			return
		}
	}

	tx.From = msg.DelegatorAddr
	tx.To = msg.ValidatorAddr
	tx.Memo = srcTx.Tx.Body.Memo

	token := p.normalizeToken(msg.Amount.Denom)

	switch msg.Type {
	case MsgDelegate:
		tx.Type = types.TxStakeDelegate
	case MsgUndelegate:
		tx.Type = types.TxStakeUndelegate
	case MsgWithdrawDelegatorReward:
		tx.Type = types.TxStakeClaimRewards
		amount = srcTx.Logs.GetWithdrawRewardValue(p.Denom)
	}

	tx.Metadata = &types.Transfer{
		Asset:  p.Coin().TokenAssetID(token),
		Amount: types.Amount(amount),
	}
}

func (p *Platform) normalizeToken(token string) string {
	if p.isCoin(token) {
		token = ""
	}

	return token
}

func (p *Platform) isCoin(asset string) bool {
	return asset == string(DenomAtom) ||
		asset == string(DenomOsmosis) ||
		asset == string(DenomKava) ||
		asset == string(DenomLuna)
}
