package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kingsleyocran/simple_bank_bankend/db/sqlc"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

//getTransferRequest
//Takes an id as a URI parameter Eg. accounts/:id
type getTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

//getTransfer request and response handler function
func (server *Server) getTransfer(ctx *gin.Context) {
	//Note that for URI parameters we use ShouldBindUri
	var req getTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.store.GetTransfer(ctx, req.ID)
	if err != nil {
		//if returns emptyRow error: sql.ErrNoRows
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

//listTransferRequest struct to get paginated data
//Take a query string instead. We use `form:"variable"`
type listTransferRequest struct {
	FromAccountID int64 `form:"from_account_id" binding:"required"`
	ToAccountID   int64 `form:"to_account_id" binding:"required"`
	PageID        int32 `form:"page_id" binding:"required,min=1"`
	PageSize      int32 `form:"page_size" binding:"required,min=5,max=10"`
}

//listTransfers request and response handler function
func (server *Server) listTransfers(ctx *gin.Context) {
	var req listTransferRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Note that Limit is simply the req.PageSize.
	//Offset is the number of records that the database should skip,
	//we we have to calculate it from the page id and page size using this formula: (req.PageID - 1) * req.PageSize
	arg := db.ListTransfersParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Limit:         req.PageSize,
		Offset:        (req.PageID - 1) * req.PageSize,
	}

	transfers, err := server.store.ListTransfers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}

//function to check if account exists in our database
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
