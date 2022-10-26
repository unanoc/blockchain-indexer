package near

import "github.com/unanoc/blockchain-indexer/pkg/client"

type Client struct {
	client.Request
}

func (c *Client) GetCurrentBlockNumber() (int64, error) {
	var block Block

	if err := c.RPCCall(&block, "block", map[string]string{"finality": "final"}); err != nil {
		return 0, err
	}

	return int64(block.Header.Height), nil
}

func (c *Client) GetBlockByNumber(num int64) (ChunkDetail, error) {
	var block Block

	if err := c.RPCCall(&block, "block", map[string]int64{"block_id": num}); err != nil || len(block.Chunks) == 0 {
		return ChunkDetail{}, err
	}

	var chunk ChunkDetail
	if err := c.RPCCall(&chunk, "chunk", []string{block.Chunks[0].Hash}); err != nil {
		return ChunkDetail{}, err
	}

	chunk.Header.Timestamp = block.Header.Timestamp

	return chunk, nil
}

func (c *Client) GetVersion() (string, error) {
	var nodeStatus NodeStatus

	if err := c.RPCCall(&nodeStatus, "status", nil); err != nil {
		return "", err
	}

	return nodeStatus.Version.Version, nil
}
