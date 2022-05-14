package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kingsleyocran/simple_bank_bankend/db/sqlc"
)

/*
Define a new Server struct. This Server will serves all HTTP requests for our
banking service. It will have 2 fields:

@db.Store: It will allow us to interact with the database when processing API requests from clients.
@gin.Engine. This router will help us send each API request to the correct handler for processing.
*/

type Server struct {
	store  db.Store
	router *gin.Engine
}

/*
NewServer takes a db.Store as input, and return a Server. This function will create a new Server
instance, and setup all HTTP API routes for our service on that server.
*/

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	//router.POST("/entries", server.createEntry)
	router.GET("/entries/:id", server.getEntry)
	router.GET("/entries", server.listEntries)

	router.POST("/transfers", server.createTransfer)
	router.GET("/transfers/:id", server.getTransfer)
	router.GET("/transfers", server.listTransfers)

	server.router = router
	return server
}

//Starts a server on an address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

//Custom error function to handle errors
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
