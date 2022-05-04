package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kingsleyocran/simple_bank_bankend/db/sqlc"
)

//Struct to handle request from http
//When a new account is created, its initial balance should always be 0, so we can remove the balance field.
//We also use binding for validation
//To be safe copy the createAccountParams to maintain the names and json
type createAccountRequest struct {
	OwnerName string `json:"owner_name" binding:"required"`
	Currency  string `json:"currency" binding:"required,oneof=USD EUR"`
}

//createAccount request and response handler function
func (server *Server) createAccount(ctx *gin.Context) {
	//create request and validate
	//Note that we used ShouldBindJSON for json parameters
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//create arg for query
	arg := db.CreateAccountParams{
		OwnerName: req.OwnerName,
		Currency:  req.Currency,
		Balance:   0,
	}

	//run createAccount with arg
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		//return error as response
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//return StatusOk if passed as response
	ctx.JSON(http.StatusOK, account)
}

//Get account request
//Takes an id as a URI parameter Eg. accounts/:id
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

//createAccount request and response handler function
func (server *Server) getAccount(ctx *gin.Context) {
	//Note that for URI parameters we use ShouldBindUri
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		//if returns emptyRow error: sql.ErrNoRows
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

//listAccount struct to get paginated data
//Take a query string instead. We use `form:"variable"`
type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

//listAccount request and response handler function
func (server *Server) listAccount(ctx *gin.Context) {
	//Note that we used ShouldBindQuery for quesry string validation
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//Note that Limit is simply the req.PageSize.
	//Offset is the number of records that the database should skip,
	//we we have to calculate it from the page id and page size using this formula: (req.PageID - 1) * req.PageSize
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
