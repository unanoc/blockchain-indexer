package transactions

// Transactions
type (
	TxsResp struct {
		StatusCode int   `json:"status_code"`
		TotalCount int64 `json:"total_count"`
		TotalPages int   `json:"total_pages"`
		PageNumber int   `json:"page_number"`
		Limit      int   `json:"limit"`
		Txs        []Tx  `json:"txs"`
	}

	TxResp struct {
		StatusCode int `json:"status_code"`
		Tx         Tx  `json:"tx"`
	}

	Tx struct {
		Hash      string      `json:"hash"`
		Chain     string      `json:"chain"`
		Height    uint64      `json:"height"`
		From      string      `json:"from"`
		To        string      `json:"to"`
		Status    string      `json:"status"`
		Type      string      `json:"type"`
		Sequence  uint64      `json:"sequence"`
		Fee       string      `json:"fee"`
		Data      interface{} `json:"data"`
		Timestamp int64       `json:"timestamp"`
	}
)
