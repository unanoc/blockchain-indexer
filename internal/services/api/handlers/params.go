package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/unanoc/blockchain-indexer/internal/services/api/httperr"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
)

const (
	ErrInvalidAddress    = "invalid address"
	ErrInvalidPage       = "invalid page"
	ErrInvalidLimit      = "invalid limit"
	ErrInvalidChain      = "invalid chain"
	ErrChainDoesNotExist = "chain doesn't exist"
	ErrEmptyHash         = "empty hash"

	defaultPage  = 1
	defaultLimit = 30
)

type txsQueryParams struct {
	Chain  string
	Page   int
	Limit  int
	Recent bool
}

func validateTransactionsParams(c *gin.Context) (*txsQueryParams, *httperr.Error) {
	chain, ok := c.GetQuery("chain")
	if !ok {
		return nil, httperr.NewError(http.StatusBadRequest, ErrInvalidChain)
	}

	_, ok = coin.Chains[chain]
	if !ok {
		return nil, httperr.NewError(http.StatusBadRequest, ErrChainDoesNotExist)
	}

	pageStr, ok := c.GetQuery("page")
	if !ok {
		pageStr = fmt.Sprintf("%d", defaultPage)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, httperr.NewError(http.StatusBadRequest, ErrInvalidPage)
	}

	if page == 0 {
		return nil, httperr.NewError(http.StatusBadRequest, ErrInvalidPage)
	}

	limitStr, ok := c.GetQuery("limit")
	if !ok {
		limitStr = fmt.Sprintf("%d", defaultLimit)
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, httperr.NewError(http.StatusBadRequest, ErrInvalidLimit)
	}

	recent := false
	recentStr, ok := c.GetQuery("recent")
	if ok && recentStr == "true" {
		recent = true
	}

	return &txsQueryParams{
		Chain:  chain,
		Page:   page,
		Limit:  limit,
		Recent: recent,
	}, nil
}

type txQueryParams struct {
	Chain string
	Hash  string
}

func validateTransactionParams(c *gin.Context) (*txQueryParams, *httperr.Error) {
	chain, ok := c.GetQuery("chain")
	if !ok {
		return nil, httperr.NewError(http.StatusBadRequest, ErrInvalidChain)
	}

	hash := c.Param("hash")
	if hash == "" {
		return nil, httperr.NewError(http.StatusBadRequest, ErrEmptyHash)
	}

	return &txQueryParams{
		Chain: chain,
		Hash:  hash,
	}, nil
}

type txsUserQueryParams struct {
	txsQueryParams
	Address string
}

func validateUserTransactionsParams(c *gin.Context) (*txsUserQueryParams, *httperr.Error) {
	txsParams, err := validateTransactionsParams(c)
	if err != nil {
		return nil, err
	}

	address, ok := c.GetQuery("address")
	if !ok {
		return nil, httperr.NewError(http.StatusBadRequest, ErrInvalidAddress)
	}

	return &txsUserQueryParams{
		Address:        address,
		txsQueryParams: *txsParams,
	}, nil
}
