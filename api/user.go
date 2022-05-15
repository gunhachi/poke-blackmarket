package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
)

// createUserRequest represent incoming data request to be pass for user accout
type createUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	UserRole string `json:"user_role" binding:"required,oneof=LEAD GRUNT"`
}

// createUser to
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

// getUserRequest represent incoming data request to be pass for user accout
type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getUser
func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}
