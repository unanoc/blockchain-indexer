package main

import (
	"context"

	"github.com/unanoc/blockchain-indexer/internal/services/transactionconsumer"
)

func main() {
	transactionconsumer.NewApp().Run(context.Background())
}
