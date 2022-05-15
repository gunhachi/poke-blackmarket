package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
)

type Server struct {
	store *db.Store
	route *gin.Engine
}

// NewServer creates a new HTTP server and setup routes
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/user", server.createUser)

	server.route = router

	return server
}

// Start run the http server
func (server *Server) Start(address string) error {
	return server.route.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
