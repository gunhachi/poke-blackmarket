package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/gunhachi/poke-blackmarket/db/sqlc"
	"github.com/gunhachi/poke-blackmarket/token"
)

// createOrderRequest represent request payload for creating order
type createOrderRequest struct {
	UserID    int64 `json:"user_id" binding:"required"`
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int32 `json:"quantity" binding:"required"`
}

// createOrder handler for creating order data and put the transaction into database layer
func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, valid := server.validUser(ctx, req.UserID, "GRUNT")
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.UserName != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.OrderTxParams{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	order, err := server.store.OrderTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)

}

// getOrderRequest represend id of order data for binding parameter to getOrder handler
type getOrderRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getOrderUserIDReq represent id of user data for binding parameter to getOrder handler
type getOrderUserIDReq struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// getOrder handler of get order data based on given order id and responding user id
func (server *Server) getOrder(ctx *gin.Context) {
	var req getOrderRequest
	var orderID getOrderUserIDReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&orderID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, valid := server.validUser(ctx, orderID.UserID, "GRUNT")
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.UserName != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	order, err := server.store.GetPokemonOrderData(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)

}

// listOrderRequest represent parameter to list the order data
type listOrderRequest struct {
	UserID   int64 `form:"user_id"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listOrder handler to list the order on database
func (server *Server) listOrder(ctx *gin.Context) {
	var req listOrderRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, valid := server.validUser(ctx, req.UserID, "LEAD")
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.UserName != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.ListPokemonOrderDataParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	orders, err := server.store.ListPokemonOrderData(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)

}

// listOrder handler to list the order on database
func (server *Server) listOrderDetailed(ctx *gin.Context) {
	var req listOrderRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, valid := server.validUser(ctx, req.UserID, "LEAD")
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.UserName != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.ListOrderDetailedDataParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	orders, err := server.store.ListOrderDetailedData(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)

}

// cancelOrder handler to cancel pokemon order data based on id
func (server *Server) cancelOrder(ctx *gin.Context) {
	var req getOrderRequest
	var orderID getOrderUserIDReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&orderID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, valid := server.validUser(ctx, orderID.UserID, "GRUNT")
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.UserName != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, err := server.store.CancelOrderTx(ctx, db.CancelOrderParam{ID: req.ID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"deleted": req.ID})

}

// validUser check whether user data valid based on param of id,user_name, and role
func (server *Server) validUser(ctx *gin.Context, userID int64, role string) (db.User, bool) {
	user, err := server.store.GetUserAccount(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return user, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return user, false
	}

	if user.UserRole != role {
		err := fmt.Errorf("user [%d] role mismatch: %s vs %s", user.ID, user.UserRole, role)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return user, false
	}

	return user, true
}
