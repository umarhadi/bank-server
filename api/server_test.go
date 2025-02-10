package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/umarhadi/bank-server/db/sqlc"
	"github.com/umarhadi/bank-server/util"
)

func TestNewServer(t *testing.T) {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	testCases := []struct {
		name          string
		config        util.Config
		checkResponse func(t *testing.T, server *Server, err error)
	}{
		{
			name:   "success",
			config: config,
			checkResponse: func(t *testing.T, server *Server, err error) {
				require.NoError(t, err)
				require.NotNil(t, server)
				require.NotNil(t, server.tokenMaker)
				require.NotNil(t, server.router)
			},
		},
		{
			name: "invalid token symmetric key",
			config: util.Config{
				TokenSymmetricKey:   "", // Invalid key
				AccessTokenDuration: time.Minute,
			},
			checkResponse: func(t *testing.T, server *Server, err error) {
				require.Error(t, err)
				require.Nil(t, server)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := db.NewStore(nil) // Using nil for now since we don't need actual DB
			server, err := NewServer(tc.config, store)
			tc.checkResponse(t, server, err)
		})
	}
}

func TestStartServer(t *testing.T) {
	store := db.NewStore(nil) // Using nil for now since we don't need actual DB
	server := newTestServer(t, store)

	// Test starting server on an invalid address
	err := server.Start("invalid:address")
	require.Error(t, err)
}
