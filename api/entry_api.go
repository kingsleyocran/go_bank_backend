package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kingsleyocran/simple_bank_bankend/db/sqlc"
)

type createEntryRequest struct {
	AccountID int64 `json:"account_id" binding:"required"`
	Amount    int64 `json:"amount" binding:"required"`
}

//CreateEntry request and response handler function
func (server *Server) createEntry(ctx *gin.Context) {
	//create request and validate
	//Note that we used ShouldBindJSON for json parameters
	var req createEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//create arg for query
	arg := db.CreateEntryParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	}

	//run CreateEntry with arg
	entry, err := server.store.CreateEntry(ctx, arg)
	if err != nil {
		//return error as response
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//return StatusOk if passed as response
	ctx.JSON(http.StatusOK, entry)
}

//Get entry request
//Takes an id as a URI parameter Eg. accounts/:id
type getEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

//getEntry request and response handler function
func (server *Server) getEntry(ctx *gin.Context) {
	//Note that for URI parameters we use ShouldBindUri
	var req getEntryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := server.store.GetEntry(ctx, req.ID)
	if err != nil {
		//if returns emptyRow error: sql.ErrNoRows
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

//listAccount struct to get paginated data
//Take a query string instead. We use `form:"variable"`
type listEntriesRequest struct {
	AccountID int64 `form:"account_id" binding:"required"`
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=10"`
}

//listEntry request and response handler function
func (server *Server) listEntries(ctx *gin.Context) {
	var req listEntriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Note that Limit is simply the req.PageSize.
	//Offset is the number of records that the database should skip,
	//we we have to calculate it from the page id and page size using this formula: (req.PageID - 1) * req.PageSize
	arg := db.ListEntriesParams{
		AccountID: req.AccountID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	entries, err := server.store.ListEntries(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}
