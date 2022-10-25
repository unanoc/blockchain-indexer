package bitcoin

import (
	"errors"
	"regexp"

	"github.com/unanoc/blockchain-indexer/pkg/client"
)

var ErrNotInSync = errors.New("not in sync to get current block number")

type Client struct {
	client.Request
}

func (c *Client) GetCurrentBlockNumber() (int64, error) {
	var nodeInfo NodeInfo
	if err := c.Get(&nodeInfo, "/api/v2", nil); err != nil {
		return 0, err
	}

	// If not in sync, latest block might not be available yet.
	if !nodeInfo.Blockbook.InSync {
		return 0, ErrNotInSync
	}

	return nodeInfo.Blockbook.BestHeight, nil
}

func (c *Client) GetVersion() (string, error) {
	var nodeInfo NodeInfo
	var version string

	if err := c.Get(&nodeInfo, "/api/v2", nil); err != nil {
		return version, err
	}

	if nodeInfo.Backend == nil || nodeInfo.Backend.Subversion == "" {
		return version, nil
	}

	re := regexp.MustCompile(`:([0-9\.-_]*).*\/`) // e.g. /MagicBean:4.6.0-2/
	matches := re.FindStringSubmatch(nodeInfo.Backend.Subversion)

	if len(matches) == 2 {
		version = matches[1]
	}

	return version, nil
}
