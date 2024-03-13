package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	db "gopsql/banking/db/sqlc"
	"gopsql/banking/db/sqlc/mocks"
	"gopsql/banking/token"
	tokenMocks "gopsql/banking/token/mocks"
	"gopsql/banking/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func randomAccount() db.Account {
	return db.Account{
		ID:        1,
		Owner:     "John",
		Balance:   20,
		Currency:  db.CurrencyEUR,
		CreatedAt: time.Now(),
	}
}

func TestCreateAccount(t *testing.T) {
	mockStore := mocks.NewStore(t)
	expectedAccount := randomAccount()
	mockStore.EXPECT().CreateAccount(
		mock.AnythingOfType("*gin.Context"),
		db.CreateAccountParams{
			Owner:    expectedAccount.Owner,
			Currency: expectedAccount.Currency,
			Balance:  0,
		},
	).Return(expectedAccount, nil).Once()

	config := util.LoadConfig()
	tokenMaker, err := token.NewJWTMaker(config.PrivateKey, config.PublicKey)
	require.NoError(t, err)
	server := NewServer(mockStore, tokenMaker)

	w := httptest.NewRecorder()
	jsonValue, _ := json.Marshal(expectedAccount)
	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)

	server.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var receivedAccount db.Account
	err = json.Unmarshal((w.Body.Bytes()), &receivedAccount)

	require.NoError(t, err)
	require.EqualExportedValues(t, expectedAccount, receivedAccount)
}

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(mockStore *mocks.Store) {
				mockStore.EXPECT().GetAccount(
					mock.AnythingOfType("*gin.Context"),
					account.ID,
				).Return(account, nil).Once()
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var receivedAccount db.Account
				err := json.Unmarshal((recorder.Body.Bytes()), &receivedAccount)
				require.NoError(t, err)
				require.EqualExportedValues(t, account, receivedAccount)
			},
		},
		{
			name:      "NotFound",
			accountID: int64(444),
			buildStubs: func(mockStore *mocks.Store) {
				mockStore.EXPECT().GetAccount(
					mock.AnythingOfType("*gin.Context"),
					int64(444),
				).Return(db.Account{}, sql.ErrNoRows).Once()
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: int64(555),
			buildStubs: func(mockStore *mocks.Store) {
				mockStore.EXPECT().GetAccount(
					mock.AnythingOfType("*gin.Context"),
					int64(555),
				).Return(db.Account{}, sql.ErrConnDone).Once()
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			buildStubs: func(mockStore *mocks.Store) {
				mockStore.EXPECT().GetAccount(
					mock.Anything,
					mock.Anything,
				).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	mockStore := mocks.NewStore(t)
	tokenMaker := tokenMocks.NewTokenMaker(t)
	for i := range testCases {
		tc := testCases[i]
		t.Run(
			tc.name,
			func(t *testing.T) {
				tc.buildStubs(mockStore)
				server := NewServer(mockStore, tokenMaker)
				w := httptest.NewRecorder()

				req, err := http.NewRequest(
					"GET",
					fmt.Sprintf("/accounts/%d", tc.accountID),
					nil,
				)

				require.NoError(t, err)

				server.router.ServeHTTP(w, req)

				tc.checkResponse(t, w)
			},
		)
	}
}
