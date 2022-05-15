package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
)

// Server server http request of pokemon blakmarket services
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer create new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/user", server.createUser)
	router.GET("/user/:id", server.getUser)

	server.router = router

	return server
}

// Start run the HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
