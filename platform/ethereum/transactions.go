package ethereum

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/pkg/primitives/address"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

const hashEventTypeTransfer = "0xddf252ad"

func (p *Platform) NormalizeRawBlock(rawBlock []byte) (types.Txs, error) {
	var block Block
	if err := json.Unmarshal(rawBlock, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return p.NormalizeBlock(block, block.TxReceipts.Map()), nil
}

func (p *Platform) NormalizeBlock(block Block, receipts map[string]TransactionReceipt) types.Txs {
	var txs types.Txs

	for _, tx := range block.Transactions {
		if tx.BlockNumber == nil { // pending transaction
			continue
		}

		receipt, ok := receipts[tx.Hash]
		if !ok {
			continue
		}

		if receipt.Status == nil {
			log.WithField("chain", p.Coin().Handle).Errorf("Empty tx receipt status for tx: %s", tx.Hash)

			continue
		}

		normalized := p.NormalizeTransaction(receipt, tx, block)
		if normalized == nil {
			continue
		}

		status := types.StatusSuccess
		if (*big.Int)(receipt.Status).Cmp(big.NewInt(0)) == 0 {
			status = types.StatusError
		}
		normalized.Status = status

		txs = append(txs, *normalized)
	}

	return txs
}

func (p *Platform) NormalizeTransaction(receipt TransactionReceipt, tx Transaction, block Block) *types.Tx {
	ts := (*big.Int)(block.Timestamp).Int64()

	if tx.Input == "0x" { // coin transfer
		return p.NormalizeCoinTransfer(tx, receipt, ts, block.BaseFeePerGas)
	}

	isTransfer, transferLog := p.checkTransfer(receipt.Logs)
	if isTransfer {
		return p.NormalizeTokenTransfer(tx, receipt, ts, transferLog, block.BaseFeePerGas)
	}

	return p.NormalizeContractCall(tx, receipt, ts, block.BaseFeePerGas)
}

func (p *Platform) normalizeBaseOfTx(srcTx Transaction,
	receipt TransactionReceipt, timestamp int64, baseFeePerGas *types.HexNumber,
) *types.Tx {
	l := log.WithFields(log.Fields{"chain": p.Coin().Handle, "tx": srcTx.Hash})

	addressFrom, err := address.EIP55Checksum(srcTx.From)
	if err != nil {
		l.WithError(err).Error("Error while getting eip55 checksum of ethereum address 'from'")

		return nil
	}

	addressTo, err := address.EIP55Checksum(srcTx.To)
	if err != nil {
		l.WithError(err).Error("Error while getting eip55 checksum of ethereum address 'to'")

		return nil
	}

	tx := &types.Tx{
		Hash:           srcTx.Hash,
		Chain:          p.Coin().Handle,
		From:           addressFrom,
		To:             addressTo,
		BlockCreatedAt: timestamp,
		Block:          (*big.Int)(srcTx.BlockNumber).Uint64(),
		Fee: types.Fee{
			Asset: p.Coin().AssetID(),
		},
	}

	if srcTx.Nonce == nil {
		l.Errorf("nonce is nil")

		return nil
	}
	tx.Sequence = (*big.Int)(srcTx.Nonce).Uint64()

	feeValue, err := srcTx.Fee((*big.Int)(baseFeePerGas), (*big.Int)(receipt.GasUsed))
	if err != nil {
		l.WithError(err).Error("cannot compute fee")

		return nil
	}

	tx.Fee.Amount = types.Amount(feeValue)

	return tx
}

func (p *Platform) NormalizeCoinTransfer(srcTx Transaction,
	receipt TransactionReceipt, timestamp int64, baseFeePerGas *types.HexNumber,
) *types.Tx {
	tx := p.normalizeBaseOfTx(srcTx, receipt, timestamp, baseFeePerGas)
	if tx == nil {
		return nil
	}
	tx.Type = types.TxTransfer

	if srcTx.Value == nil {
		log.WithField("tx", srcTx.Hash).Error("transfer value is nil")

		return nil
	}
	tx.Metadata = &types.Transfer{
		Asset:  p.Coin().AssetID(),
		Amount: types.Amount((*big.Int)(srcTx.Value).String()),
	}

	return tx
}

func (p *Platform) NormalizeTokenTransfer(srcTx Transaction, receipt TransactionReceipt,
	timestamp int64, eventLog EventLog, baseFeePerGas *types.HexNumber,
) *types.Tx {
	tx := p.normalizeBaseOfTx(srcTx, receipt, timestamp, baseFeePerGas)
	if tx == nil {
		return nil
	}

	l := log.WithFields(log.Fields{"chain": p.Coin().Handle, "tx": tx.Hash})

	addressFrom, err := address.EIP55Checksum(common.HexToAddress(eventLog.Topics[1]).Hex())
	if err != nil {
		l.WithError(err).Errorf("Could not get eip55 checksum of ethereum token sender")

		return nil
	}

	addressTo, err := address.EIP55Checksum(common.HexToAddress(eventLog.Topics[2]).Hex())
	if err != nil {
		l.WithError(err).Errorf("Could not get eip55 checksum of ethereum token recipient")

		return nil
	}

	value, err := hexDataToBigInt(eventLog.Data)
	if err != nil {
		l.WithError(err).Error("Hex conversion to big int error")

		return nil
	}

	token, err := address.EIP55Checksum(eventLog.Address) // Use address of a contract that was transferred
	if err != nil {
		l.WithError(err).Error("EIP55 checksum getting error")

		return nil
	}

	tx.From = addressFrom
	tx.To = addressTo
	tx.Type = types.TxTransfer
	tx.Metadata = &types.Transfer{
		Asset:  p.Coin().TokenAssetID(token), // Token contract
		Amount: types.Amount(value.String()),
	}

	return tx
}

func (p *Platform) NormalizeContractCall(srcTx Transaction,
	receipt TransactionReceipt, timestamp int64, baseFeePerGas *types.HexNumber,
) *types.Tx {
	tx := p.normalizeBaseOfTx(srcTx, receipt, timestamp, baseFeePerGas)
	if tx == nil {
		return nil
	}

	tx.Type = types.TxContractCall

	if srcTx.Value == nil {
		log.WithField("tx", srcTx.Hash).Error("contract call value is nil")

		return nil
	}
	tx.Metadata = &types.ContractCall{
		Asset:  p.Coin().AssetID(),
		Amount: types.Amount((*big.Int)(srcTx.Value).String()),
		Input:  srcTx.Input,
	}

	return tx
}

//nolint:goerr113
func hexDataToBigInt(data string) (*big.Int, error) {
	data = strings.Replace(data, "0x", "", 1)

	value, ok := new(big.Int).SetString(data, 16)
	if !ok {
		return nil, fmt.Errorf("cannot convert data to *big.Int: %s", data)
	}

	return value, nil
}

func (p *Platform) checkTransfer(logs []EventLog) (bool, EventLog) {
	var transferLog EventLog
	transferLogsNum := 0
	for i, eventLog := range logs {
		topics := eventLog.Topics
		if len(topics) == 3 && strings.HasPrefix(topics[0], hashEventTypeTransfer) {
			transferLog = logs[i]
			transferLogsNum++
		}
	}

	return transferLogsNum == 1, transferLog
}
