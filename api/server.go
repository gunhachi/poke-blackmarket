package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
	"github.com/gunhachi/poke-blackmarket/token"
	"github.com/gunhachi/poke-blackmarket/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	route      *gin.Engine
}

// NewServer creates a new HTTP server and setup routes
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/account", server.createAccountLog)
	router.POST("/account/login", server.loginAccount)
	router.GET("/pokemon-api/:name", server.getDataPokemonApi)

	authRoute := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoute.GET("/user", server.listUser)
	authRoute.GET("/user/:id", server.getUser)
	authRoute.POST("/user", server.createUser)

	authRoute.POST("/pokemon", server.createPokemon)
	authRoute.GET("/pokemon", server.listPokemon)
	authRoute.GET("/pokemon/:id", server.getPokemon)
	authRoute.PUT("/pokemon/:id", server.updatePokemon)

	authRoute.POST("/order", server.createOrder)
	authRoute.GET("/order/:id", server.getOrder)
	authRoute.DELETE("/order/:id", server.cancelOrder)

	authRoute.GET("/order", server.listOrder)
	authRoute.GET("/order-detailed", server.listOrderDetailed)
	authRoute.PUT("/user/:id", server.updateUser)

	server.route = router

}

// Start run the http server
func (server *Server) Start(address string) error {
	return server.route.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
