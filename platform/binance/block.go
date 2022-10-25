package binance

import (
	"encoding/json"
	"fmt"
)

func (p *Platform) GetCurrentBlockNumber() (int64, error) {
	block, err := p.dexClient.GetCurrentBlockNumber()
	if err != nil {
		return 0, err
	}

	return block, nil
}

func (p *Platform) GetBlockByNumber(num int64) ([]byte, error) {
	block, err := p.client.GetBlockByNumber(num)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(block)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block: %w", err)
	}

	return data, nil
}
