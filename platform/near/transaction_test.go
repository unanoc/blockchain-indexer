package near

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unanoc/blockchain-indexer/pkg/mock"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
)

func TestNormalizeChunk(t *testing.T) {
	var chunk ChunkDetail
	err := mock.JSONModelFromFilePath("mocks/tx_input.json", &chunk)
	assert.NoError(t, err)

	p := Init(coin.NEAR, "")

	result := p.NormalizeBlock(chunk)

	resultJSON, err := json.Marshal(result)
	assert.NoError(t, err)

	expectedJSON, err := mock.JSONStringFromFilePath("mocks/tx_expected_output.json")
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(resultJSON))
}
