package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
	"github.com/mtslzr/pokeapi-go"
)

type createPokemonRequest struct {
	PokeName  string `json:"poke_name" binding:"required"`
	Status    string `json:"status" binding:"required"`
	PokePrice int64  `json:"poke_price" binding:"required"`
	PokeStock int64  `json:"poke_stock" binding:"required"`
}

func (server *Server) createPokemon(ctx *gin.Context) {
	var req createPokemonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreatePokemonDataParams{
		PokeName:  req.PokeName,
		Status:    req.Status,
		PokePrice: req.PokePrice,
		PokeStock: req.PokeStock,
	}

	poke, err := server.store.CreatePokemonData(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, poke)

}

type getPokemonRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getPokemon(ctx *gin.Context) {
	var req getPokemonRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	poke, err := server.store.GetPokemonData(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, poke)

}

type pokemonAPIRequest struct {
	Name string `uri:"name" binding:"required"`
}

type pokemonAPIResponse struct {
	Name    string `json:"name"`
	BaseExp int    `json:"base_experience"`
}

func (server *Server) getDataPokemonApi(ctx *gin.Context) {
	var req pokemonAPIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	poke, err := pokeapi.Pokemon(req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := pokemonAPIResponse{
		Name:    poke.Name,
		BaseExp: poke.BaseExperience,
	}

	ctx.JSON(http.StatusOK, resp)

}

type listPokeRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listPokemon(ctx *gin.Context) {
	var req listPokeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListPokemonDataParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	pokes, err := server.store.ListPokemonData(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, pokes)

}

type updatePokemonData struct {
	// ID        int64  `json:"id" binding:"required,min=1"`
	Status    string `json:"status"`
	PokePrice int64  `json:"poke_price"`
	PokeStock int64  `json:"poke_stock"`
}

func (server *Server) updatePokemon(ctx *gin.Context) {
	var req getPokemonRequest
	var dataReq updatePokemonData
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&dataReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdatePokemonDataParams{
		ID:        req.ID,
		Status:    dataReq.Status,
		PokePrice: dataReq.PokePrice,
		PokeStock: dataReq.PokeStock,
	}

	poke, err := server.store.UpdatePokemonData(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, poke)

}
