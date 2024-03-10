package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateAndVerifyJWTToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker()
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	username := "myUsername"
	token, err := jwtMaker.CreateToken(username, 2*time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	paylaod, err := jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, paylaod)
	require.Equal(t, paylaod.Username, username)

	// Test token expiration
	time.Sleep(3 * time.Second)
	_, err = jwtMaker.VerifyToken(token)
	require.Error(t, err)
}
