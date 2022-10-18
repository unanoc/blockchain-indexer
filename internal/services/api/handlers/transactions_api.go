package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/controllers/transactions"
)

type TransactionsAPI struct {
	controller *transactions.Controller
}

func NewTransactionsAPI(db repository.Storage) API {
	return &TransactionsAPI{
		controller: transactions.NewController(db),
	}
}

// GetTransactions godoc
// @Description  Returns all transaction list by creation date order(asc/desc)
// @Tags         Transactions
// @Produce      json
// @Param        chain query string true "Chain"
// @Param        page query int false "Page for pagination"
// @Param        limit query int false "The limit of the number of items"
// @Param        recent query bool false "Enable desc order"
// @Success      200  {object}  transactions.TxsResp
// @Failure      400  {object}  httperr.Error
// @Failure      500  {object}  httperr.Error
// @Router       /api/v1/transactions [get]
func (api *TransactionsAPI) GetTransactions(c *gin.Context) {
	params, err := validateTransactionsParams(c)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	response, err := api.controller.GetTransactions(c.Request.Context(),
		params.Chain, params.Page, params.Limit, params.Recent)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	c.JSON(response.StatusCode, response)
}

// GetTransactionByHash godoc
// @Description  Returns transaction details by hash
// @Tags         Transactions
// @Produce      json
// @Param        chain query string true "Chain"
// @Param        hash  path string true "Transaction hash"
// @Success      200  {object}  transactions.TxResp
// @Failure      400  {object}  httperr.Error
// @Failure      404  {object}  httperr.Error
// @Failure      500  {object}  httperr.Error
// @Router       /api/v1/transactions/{hash} [get]
func (api *TransactionsAPI) GetTransactionByHash(c *gin.Context) {
	params, err := validateTransactionParams(c)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	response, err := api.controller.GetTransactionByHash(c.Request.Context(), params.Chain, params.Hash)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	c.JSON(response.StatusCode, response)
}

// GetUserTransactionsByUser godoc
// @Description  Returns all user transaction list by creation date order(asc/desc)
// @Tags         Transactions
// @Produce      json
// @Param        chain query string true "Chain"
// @Param        address query string true "Address that made transactions"
// @Param        page query int false "Page for pagination"
// @Param        limit query int false "The limit of the number of items"
// @Param        recent query bool false "Enable desc order"
// @Success      200  {object}  transactions.TxsResp
// @Failure      400  {object}  httperr.Error
// @Failure      500  {object}  httperr.Error
// @Router       /api/v1/transactions/user [get]
func (api *TransactionsAPI) GetTransactionsByUser(c *gin.Context) {
	params, err := validateUserTransactionsParams(c)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	response, err := api.controller.GetTransactionsByUser(c.Request.Context(),
		params.Chain, params.Address, params.Page, params.Limit, params.Recent)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	c.JSON(response.StatusCode, response)
}
