package types

import "github.com/unanoc/blockchain-indexer/pkg/primitives/coin"

type (
	Amount          string
	Asset           string
	Status          string
	TransactionType string
)

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"

	TxTransfer          TransactionType = "transfer"
	TxContractCall      TransactionType = "contract_call"
	TxStakeClaimRewards TransactionType = "stake_claim_rewards"
	TxStakeDelegate     TransactionType = "stake_delegate"
	TxStakeUndelegate   TransactionType = "stake_undelegate"
	TxStakeRedelegate   TransactionType = "stake_redelegate"
)

type (
	Block struct {
		Number int64 `json:"number"`
		Txs    []Tx  `json:"txs"`
	}

	Txs []Tx

	Tx struct {
		Hash           string          `json:"hash" bson:"hash"`
		Chain          string          `json:"chain" bson:"chain"`
		From           string          `json:"from" bson:"from"`
		To             string          `json:"to" bson:"to"`
		BlockCreatedAt int64           `json:"block_created_at" bson:"block_created_at"`
		Block          uint64          `json:"block" bson:"block"`
		Sequence       uint64          `json:"sequence" bson:"status"`
		Status         Status          `json:"status" bson:"status"`
		Type           TransactionType `json:"type" bson:"type"`
		Memo           string          `json:"memo,omitempty"`
		Fee            Fee             `json:"fee" bson:"type"`
		Metadata       interface{}     `json:"metadata" bson:"metadata"`
	}

	Fee struct {
		Asset  coin.AssetID `json:"asset"`
		Amount Amount       `json:"amount"`
	}
)

// Tx metadata types
type (
	Transfer struct {
		Asset  coin.AssetID `json:"asset"`
		Amount Amount       `json:"amount"`
	}

	ContractCall struct {
		Asset  coin.AssetID `json:"asset"`
		Amount Amount       `json:"amount"`
		Input  string       `json:"input"`
	}
)
