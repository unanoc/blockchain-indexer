package bitcoin

import (
	"github.com/unanoc/blockchain-indexer/pkg/client"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
	"github.com/unanoc/blockchain-indexer/pkg/sentry"
)

type Platform struct {
	coin   uint
	client Client
}

func Init(coin uint, url string) *Platform {
	return &Platform{
		coin:   coin,
		client: Client{client.InitJSONClient(url, sentry.DefaultSentryErrorHandler())},
	}
}

func (p Platform) Coin() coin.Coin {
	return coin.Coins[p.coin]
}

func (p *Platform) UpdateNodeConnection(url string) {
	p.client = Client{client.InitJSONClient(url, sentry.DefaultSentryErrorHandler())}
}
