package binance

import (
	"fmt"

	"github.com/unanoc/blockchain-indexer/pkg/client"
)

type Client struct {
	client.Request
}

type DexClient struct {
	client.Request
}

func (c *DexClient) GetCurrentBlockNumber() (int64, error) {
	var result NodeInfoResponse
	if err := c.Get(&result, "/api/v1/node-info", nil); err != nil {
		return 0, err
	}

	return int64(result.SyncInfo.LatestBlockHeight), nil
}

func (c *Client) GetBlockByNumber(blockNumber int64) (Block, error) {
	var result Block
	if err := c.Get(&result, fmt.Sprintf("/api/v1/blocks/%d/txs", blockNumber), nil); err != nil {
		return Block{}, err
	}

	return result, nil
}

func (c *DexClient) GetVersion() (string, error) {
	var result NodeInfoResponse
	if err := c.Get(&result, "/api/v1/node-info", nil); err != nil {
		return "", err
	}

	return result.NodeInfo.Version, nil
}
