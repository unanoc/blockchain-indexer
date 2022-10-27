package binance

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/unanoc/blockchain-indexer/pkg/mock"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
)

func TestNormalizeTransactions(t *testing.T) {
	var block Block
	err := mock.JSONModelFromFilePath("mocks/txs_input.json", &block)
	assert.NoError(t, err)

	p := Init(coin.BINANCE, "", "")

	result := p.NormalizeBlock(block)

	resultJSON, err := json.Marshal(result)
	assert.NoError(t, err)

	expectedJSON, err := mock.JSONStringFromFilePath("mocks/txs_expected.json")
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(resultJSON))
}

func TestNormalizeTransferTransaction_CatchPanics(t *testing.T) {
	testCases := []struct {
		name    string
		tx      Tx
		assetID coin.AssetID
	}{
		{name: "tx is empty", tx: Tx{}, assetID: coin.Binance().AssetID()},
	}

	p := Init(coin.BINANCE, "", "")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() { _, _ = p.normalizeTransferTransaction(tc.tx, tc.assetID) })
		})
	}
}

func TestNormalizeDelegation_CatchPanics(t *testing.T) {
	testCases := []struct {
		name    string
		tx      Tx
		assetID coin.AssetID
	}{
		{name: "tx is empty", tx: Tx{}, assetID: coin.Binance().AssetID()},
		{name: "tx i/o are nil", tx: Tx{
			Data: `{"inputs": null, "outputs": null}`,
		}, assetID: coin.Binance().AssetID()},
		{name: "tx i/o are empty", tx: Tx{
			Data: `{"inputs": [{}], "outputs": [{}]}`,
		}, assetID: coin.Binance().AssetID()},
		{name: "tx input is empty", tx: Tx{
			Data: `{"inputs": [{}], "outputs": [{"amounts": [{"amount": 1234, "denom": "bnb"}]}]}`,
		}, assetID: coin.Binance().AssetID()},
	}

	p := Init(coin.BINANCE, "", "")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() { _, _ = p.normalizeDelegationTransaction(tc.tx, tc.assetID) })
		})
	}
}
