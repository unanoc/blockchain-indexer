package platform

import (
	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
	"github.com/unanoc/blockchain-indexer/platform/binance"
	"github.com/unanoc/blockchain-indexer/platform/cosmos"
	"github.com/unanoc/blockchain-indexer/platform/ethereum"
	"github.com/unanoc/blockchain-indexer/platform/near"
)

type (
	Platform interface {
		Coin() coin.Coin
		GetCurrentBlockNumber() (int64, error)
		GetBlockByNumber(num int64) ([]byte, error)
		GetVersion() (string, error)
		NormalizeRawBlock(rawBlock []byte) (types.Txs, error)
		UpdateNodeConnection(url string)
	}

	Platforms map[string]Platform
)

func InitPlatforms() Platforms {
	return Platforms{
		coin.Binance().Handle: binance.Init(coin.BINANCE, config.Default.Platforms.Binance.Node,
			config.Default.Platforms.Binance.Dex),
		coin.Cosmos().Handle:     cosmos.Init(coin.COSMOS, cosmos.DenomAtom, config.Default.Platforms.Cosmos.Node),
		coin.Ethereum().Handle:   ethereum.Init(coin.ETHEREUM, config.Default.Platforms.Ethereum.Node),
		coin.Smartchain().Handle: ethereum.Init(coin.SMARTCHAIN, config.Default.Platforms.Smartchain.Node),
		coin.Near().Handle:       near.Init(coin.NEAR, config.Default.Platforms.Near.Node),
	}
}

//nolint:gofumpt
func GetPlatform(chain string, url string) Platform {
	switch chain {
	case coin.Binance().Handle:
		return binance.Init(coin.BINANCE, "", url)
	case coin.Cosmos().Handle:
		return cosmos.Init(coin.COSMOS, cosmos.DenomAtom, url)
	case coin.Ethereum().Handle:
		return ethereum.Init(coin.ETHEREUM, url)
	case coin.Smartchain().Handle:
		return ethereum.Init(coin.SMARTCHAIN, url)
	case coin.Near().Handle:
		return near.Init(coin.NEAR, url)
	default:
		return nil
	}
}
