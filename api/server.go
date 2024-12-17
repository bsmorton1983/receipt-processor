package api

import (
	db "github.com/bsmorton1983/receipt_processor/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our receipt processing service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/receipts/process", server.processReceipt)
	router.GET("/receipts/:id/points", server.getPoints)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
