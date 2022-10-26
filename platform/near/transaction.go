package near

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

var ErrTransferActionUnmarshal = errors.New("unable marshaling to transfer action struct")

func (p *Platform) NormalizeRawBlock(rawBlock []byte) (types.Txs, error) {
	var block ChunkDetail
	if err := json.Unmarshal(rawBlock, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return p.NormalizeBlock(block), nil
}

func (p *Platform) NormalizeBlock(block ChunkDetail) types.Txs {
	normalized := make(types.Txs, 0)

	for _, tx := range block.Transactions {
		if len(tx.Actions) != 1 {
			continue
		}

		transfer, err := mapTransfer(tx.Actions[0])
		if err != nil {
			continue
		}

		assetID := coin.Near().AssetID()

		normalized = append(normalized, types.Tx{
			Hash:  tx.Hash,
			Chain: p.Coin().Handle,
			From:  tx.SignerID,
			To:    tx.ReceiverID,
			Fee: types.Fee{
				Asset:  assetID,
				Amount: "0",
			},
			BlockCreatedAt: int64(block.Header.Timestamp / 1_000_000_000),
			Block:          block.Header.Height,
			Status:         types.StatusSuccess,
			Sequence:       uint64(tx.Nonce),
			Type:           types.TxTransfer,
			Metadata: &types.Transfer{
				Asset:  assetID,
				Amount: types.Amount(transfer.Transfer.Deposit),
			},
		})
	}

	return normalized
}

func mapTransfer(i interface{}) (action TransferAction, err error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bytes, &action); err != nil {
		return
	}

	if action.Transfer.Deposit == "" {
		err = ErrTransferActionUnmarshal
	}

	return
}
