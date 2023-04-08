package gapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/umarhadi/bank-server/db/sqlc"
	"github.com/umarhadi/bank-server/util"
	"github.com/umarhadi/bank-server/worker"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}