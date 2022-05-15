package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
)

type createUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	UserRole string `json:"user_role" binding:"required,oneof= LEAD GRUNTs"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserAccountParams{
		UserName: req.UserName,
		UserRole: req.UserRole,
	}

	user, err := server.store.CreateUserAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}
