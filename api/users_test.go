package api

import (
	"bytes"
	"encoding/json"
	db "gopsql/banking/db/sqlc"
	"gopsql/banking/db/sqlc/mocks"
	tokenMocks "gopsql/banking/token/mocks"
	"gopsql/banking/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	mockStore := mocks.NewStore(t)
	reqParam := createUserRequest{
		Username: "mohammad",
		FullName: "mohammad Dehbashi",
		Password: "123456",
		Email:    "mohammad@m.com",
	}

	mockStore.EXPECT().CreateUser(
		mock.AnythingOfType("*gin.Context"),
		mock.MatchedBy(func(arg interface{}) bool {
			params, ok := arg.(db.CreateUserParams)
			if !ok {
				return false
			}
			err := util.CheckPassword(reqParam.Password, params.HashedPassword)
			return err == nil

		}),
	).Return(
		db.User{
			Username:          reqParam.Username,
			FullName:          reqParam.FullName,
			Email:             reqParam.Email,
			CreatedAt:         time.Now(),
			PasswordChangedAt: time.Now(),
		},
		nil,
	).Once()

	tokenMaker := tokenMocks.NewTokenMaker(t)
	server := NewServer(mockStore, tokenMaker)

	w := httptest.NewRecorder()
	jsonValue, _ := json.Marshal(reqParam)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)

	server.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var receivedAccount db.Account
	err = json.Unmarshal((w.Body.Bytes()), &receivedAccount)

	require.NoError(t, err)
	// require.EqualExportedValues(t, expectedAccount, receivedAccount)
}
