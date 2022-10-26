package near

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/unanoc/blockchain-indexer/pkg/client"
)

const errorMissingBlock = -32000

func (p *Platform) GetCurrentBlockNumber() (int64, error) {
	return p.client.GetCurrentBlockNumber()
}

func (p *Platform) GetBlockByNumber(num int64) ([]byte, error) {
	block, err := p.client.GetBlockByNumber(num)
	if err != nil { // near won't return old blocks
		rpcError := &client.RPCError{}

		if errors.As(err, &rpcError) && rpcError.Code == errorMissingBlock {
			return nil, nil
		}

		return nil, err
	}

	data, err := json.Marshal(block)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block: %w", err)
	}

	return data, nil
}
