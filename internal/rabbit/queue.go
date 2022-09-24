package rabbit

import (
	"github.com/unanoc/blockchain-indexer/pkg/mq"
)

const ExchangeTransactionsParsed mq.ExchangeName = "transactions.parsed"

const QueueTransactionsSave mq.QueueName = "transactions_save"

func NewConsumerOptions(workers int, maxRetries ...int) *mq.ConsumerOptions {
	options := mq.DefaultConsumerOptions(workers)

	if len(maxRetries) > 0 {
		options.MaxRetries = maxRetries[0]
	}

	return options
}
