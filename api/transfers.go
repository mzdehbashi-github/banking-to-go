package api

import (
	"database/sql"
	"fmt"
	db "gopsql/banking/db/sqlc"
	"gopsql/banking/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	isValid := s.validateCurrency(ctx, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}

	isValid = s.validateCurrency(ctx, req.ToAccountID, req.Currency)
	if !isValid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	res, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (s *Server) validateCurrency(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return false
	}

	if account.Currency != db.Currency(currency) {
		err := fmt.Errorf("account %d currency missmatch, expected %v, got %v", accountID, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return false
	}
	return true
}
