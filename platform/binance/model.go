package binance

const (
	Transfer   TxType = "TRANSFER"
	Delegate   TxType = "SIDECHAIN_DELEGATE"
	Undelegate TxType = "SIDECHAIN_UNDELEGATE"
)

type TxType string

type (
	NodeInfoResponse struct {
		NodeInfo struct {
			Version string `json:"version"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHeight int `json:"latest_block_height"`
		} `json:"sync_info"`
	}

	Block struct {
		Txs []Tx `json:"txs"`
	}

	Tx struct {
		TxHash      string      `json:"hash"`
		BlockHeight int         `json:"blockHeight"`
		BlockTime   uint64      `json:"blockTime"`
		TxType      TxType      `json:"type"`
		FromAddr    interface{} `json:"fromAddr"`
		ToAddr      interface{} `json:"toAddr"`
		Amount      *uint64     `json:"amount"`
		TxAsset     string      `json:"asset"`
		TxFee       uint64      `json:"fee"`
		OrderID     string      `json:"orderId,omitempty"`
		Code        int         `json:"code"`
		Data        string      `json:"data"`
		Memo        string      `json:"memo"`
		Source      int         `json:"source"`
		Sequence    int         `json:"sequence"`
	}

	DelegationData struct {
		ValidatorAddress string `json:"validatorAddr"`
		Delegation       struct {
			Amount uint64 `json:"amount"`
		} `json:"delegation"`
		Amount struct {
			Amount uint64 `json:"amount"`
		} `json:"amount"`
	}
)
