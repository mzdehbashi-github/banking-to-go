package api

import (
	"database/sql"
	db "gopsql/banking/db/sqlc"
	"gopsql/banking/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: db.Currency(req.Currency),
		Balance:  int64(0),
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, account)

}

type listAccountsRequest struct {
	PageNum  int `form:"page_num" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=10"`
}

func (s *Server) ListAccounts(ctx *gin.Context) {
	req := listAccountsRequest{
		PageNum:  1,
		PageSize: 10,
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	accounts, err := s.store.ListAccounts(
		ctx,
		db.ListAccountsParams{
			Limit:  int32(req.PageSize),
			Offset: int32((req.PageNum - 1) * req.PageSize),
		},
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
