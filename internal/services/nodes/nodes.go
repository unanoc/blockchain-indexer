package nodes

import (
	"context"
	"fmt"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
)

//nolint:gochecknoglobals
var nodesList = map[string][]models.Node{
	"binance": {
		{
			Chain:      "binance",
			Scheme:     "https",
			Host:       "dex.binance.org",
			Enabled:    true,
			Monitoring: true,
		},
	},
	"cosmos": {
		{
			Chain:      "cosmos",
			Scheme:     "https",
			Host:       "us-atom-restapi.binancechain.io",
			Enabled:    true,
			Monitoring: true,
		},
		{
			Chain:      "cosmos",
			Scheme:     "https",
			Host:       "cosmos-mainnet-rpc.allthatnode.com:1317",
			Enabled:    true,
			Monitoring: true,
		},
	},
	"ethereum": {
		{
			Chain:      "ethereum",
			Scheme:     "https",
			Host:       "ethereum-mainnet-rpc.allthatnode.com",
			Enabled:    true,
			Monitoring: true,
		},
	},
	"smartchain": {
		{
			Chain:      "smartchain",
			Scheme:     "https",
			Host:       "bsc-dataseed1.binance.org",
			Enabled:    true,
			Monitoring: true,
		},
		{
			Chain:      "smartchain",
			Scheme:     "https",
			Host:       "bsc-dataseed2.binance.org",
			Enabled:    true,
			Monitoring: true,
		},
		{
			Chain:      "smartchain",
			Scheme:     "https",
			Host:       "bsc-dataseed3.binance.org",
			Enabled:    true,
			Monitoring: true,
		},
		{
			Chain:      "smartchain",
			Scheme:     "https",
			Host:       "bsc-dataseed4.binance.org",
			Enabled:    true,
			Monitoring: true,
		},
	},
	"near": {
		{
			Chain:      "near",
			Scheme:     "https",
			Host:       "rpc.mainnet.near.org",
			Enabled:    true,
			Monitoring: true,
		},
	},
}

func AddNodesListToDB(db *postgres.Database) error {
	for _, nodesList := range nodesList {
		if err := db.InsertNodes(context.Background(), nodesList); err != nil {
			return fmt.Errorf("failed to insert nodes: %w", err)
		}
	}

	return nil
}
