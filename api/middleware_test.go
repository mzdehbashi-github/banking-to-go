package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopsql/banking/token"
	"gopsql/banking/token/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	mockTokenMaker := mocks.NewTokenMaker(t)

	testCases := []struct {
		name                  string
		authorizationHeader   string
		expectedStatusCode    int
		expectedErrorResponse string
		buildStub             func(tokenMaker *mocks.TokenMaker)
	}{
		{
			name:                  "EmptyAuthorizationHeader",
			authorizationHeader:   "",
			expectedStatusCode:    http.StatusUnauthorized,
			expectedErrorResponse: `{"error":"authorization header was not provided"}`,
			buildStub:             func(tm *mocks.TokenMaker) {},
		},
		{
			name:                  "InvalidAuthorizationHeader",
			authorizationHeader:   "invalid_header",
			expectedStatusCode:    http.StatusUnauthorized,
			expectedErrorResponse: `{"error":"invalid authorization header"}`,
			buildStub:             func(tm *mocks.TokenMaker) {},
		},
		{
			name:                  "UnsupportedAuthorizationType",
			authorizationHeader:   "basic token",
			expectedStatusCode:    http.StatusUnauthorized,
			expectedErrorResponse: `{"error":"authorization type basic is not supported"}`,
			buildStub:             func(tm *mocks.TokenMaker) {},
		},
		{
			name:                  "InvalidToken",
			authorizationHeader:   "bearer invalid_token",
			expectedStatusCode:    http.StatusUnauthorized,
			expectedErrorResponse: `{"error":"invalid token"}`,
			buildStub: func(tm *mocks.TokenMaker) {
				tm.EXPECT().VerifyToken("invalid_token").Return(nil, errors.New("invalid token")).Once()
			},
		},
		{
			name:                  "ValidToken",
			authorizationHeader:   "bearer valid_token",
			expectedStatusCode:    http.StatusOK,
			expectedErrorResponse: "",
			buildStub: func(tm *mocks.TokenMaker) {
				tm.EXPECT().VerifyToken("valid_token").Return(
					&token.Payload{},
					nil,
				).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup stubs
			tc.buildStub(mockTokenMaker)

			// Setup Gin
			r := gin.New()
			r.Use(authMiddleware(mockTokenMaker))
			r.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{})
			})

			// Create request
			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)

			req.Header.Set(authorizationHeaderKey, tc.authorizationHeader)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Check response
			require.Equal(t, tc.expectedStatusCode, w.Code)

			if tc.expectedErrorResponse != "" {
				require.JSONEq(t, tc.expectedErrorResponse, w.Body.String())
			}
		})
	}
}
