package near

// Block
type (
	Block struct {
		Header BlockHeader `json:"header"`
		Chunks []Chunk     `json:"chunks"`
	}

	BlockHeader struct {
		Height    uint64 `json:"height"`
		Timestamp uint64 `json:"timestamp"`
	}

	Chunk struct {
		Hash      string `json:"chunk_hash"`
		Height    uint64 `json:"height_created"`
		Timestamp uint64 `json:"time_stamp"`
	}

	ChunkDetail struct {
		Header       Chunk `json:"header"`
		Transactions []Tx  `json:"transactions,omitempty"`
	}

	Tx struct {
		SignerID   string        `json:"signer_id"`
		Nonce      int           `json:"nonce"`
		ReceiverID string        `json:"receiver_id"`
		Actions    []interface{} `json:"actions"`
		Hash       string        `json:"hash"`
	}

	TransferAction struct {
		Transfer Transfer `json:"Transfer"`
	}

	Transfer struct {
		Deposit string `json:"deposit"`
	}

	FunctionCall struct {
		MethodName string `json:"method_name"`
		Args       string `json:"args"`
		Gas        int64  `json:"gas"`
		Deposit    string `json:"deposit"`
	}
)

// Node version
type (
	NodeStatus struct {
		Version struct {
			Version string `json:"version"`
		} `json:"version"`
	}
)
