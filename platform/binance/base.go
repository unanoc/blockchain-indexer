package binance

import (
	"github.com/unanoc/blockchain-indexer/pkg/client"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
	"github.com/unanoc/blockchain-indexer/pkg/sentry"
)

type Platform struct {
	coin      uint
	client    Client
	dexClient DexClient
	url       string
}

func Init(coin uint, url, dex string) *Platform {
	return &Platform{
		coin:      coin,
		client:    Client{client.InitJSONClient(url, sentry.DefaultSentryErrorHandler())},
		dexClient: DexClient{client.InitJSONClient(dex, sentry.DefaultSentryErrorHandler())},
		url:       url,
	}
}

func (p Platform) Coin() coin.Coin {
	return coin.Coins[p.coin]
}

func (p *Platform) UpdateNodeConnection(url string) {
	p.client = Client{client.InitJSONClient(p.url, sentry.DefaultSentryErrorHandler())}
	p.dexClient = DexClient{client.InitJSONClient(url, sentry.DefaultSentryErrorHandler())}
}
