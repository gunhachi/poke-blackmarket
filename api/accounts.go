package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
	"github.com/gunhachi/poke-blackmarket/util"
	"github.com/lib/pq"
)

// createAccountRequest represent the param for create account
type createAccountRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

// responseAccount represent response of handler
type responseAccount struct {
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}

// buildAccountResponse build expected response
func buildAccountResponse(account db.Account) responseAccount {
	return responseAccount{
		Username:  account.Username,
		FullName:  account.FullName,
		CreatedAt: account.CreatedAt,
	}
}

// createAccountLog handler to create account
func (server *Server) createAccountLog(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountLogParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
	}

	account, err := server.store.CreateAccountLog(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := buildAccountResponse(account)

	ctx.JSON(http.StatusOK, resp)

}

// loginAccountRequest represent param for login access
type loginAccountRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

// responseLoginAccount represent response for given access such token
type responseLoginAccount struct {
	AccessToken string          `json:"access_token"`
	User        responseAccount `json:"account"`
}

// loginAccount handler to login based on existing data
func (server *Server) loginAccount(ctx *gin.Context) {
	var req loginAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccountLog(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, account.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		account.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := responseLoginAccount{
		AccessToken: accessToken,
		User:        buildAccountResponse(account),
	}

	ctx.JSON(http.StatusOK, rsp)
}
