package api

import (
	"errors"
	"gopsql/banking/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createTokenRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) CreateToken(ctx *gin.Context) {
	var req createTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("bad credentials")))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("bad credentials")))
		return
	}

	token, err := s.tokenMaker.CreateToken(req.Username, 1*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
