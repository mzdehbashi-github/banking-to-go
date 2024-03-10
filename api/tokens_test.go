package api

import (
	"bytes"
	"encoding/json"
	db "gopsql/banking/db/sqlc"
	"gopsql/banking/db/sqlc/mocks"
	"gopsql/banking/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	mockStore := mocks.NewStore(t)
	reqParam := createTokenRequest{
		Username: "mohammad",
		Password: "123456",
	}

	hasedPassword, err := util.HashPassword(reqParam.Password)
	require.NoError(t, err)

	mockStore.EXPECT().GetUser(
		mock.AnythingOfType("*gin.Context"),
		reqParam.Username,
	).Return(
		db.User{
			Username:          reqParam.Username,
			FullName:          "Mohammad",
			Email:             "m@m.com",
			CreatedAt:         time.Now(),
			PasswordChangedAt: time.Now(),
			HashedPassword:    hasedPassword,
		},
		nil,
	).Once()

	server := NewServer(mockStore)

	w := httptest.NewRecorder()
	jsonValue, _ := json.Marshal(reqParam)
	req, err := http.NewRequest("POST", "/tokens", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)

	server.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	response := struct {
		Token string `json:"token"`
	}{}
	err = json.Unmarshal((w.Body.Bytes()), &response)

	require.NoError(t, err)
	require.NotEmpty(t, response.Token)

	payload, err := server.tokenMaker.VerifyToken(response.Token)
	require.NoError(t, err)
	require.Equal(t, payload.Username, reqParam.Username)
}
