package worker

import (
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func TestNewRedisTaskDistributor(t *testing.T) {
	redisOpt := asynq.RedisClientOpt{
		Addr: "localhost:6379",
	}

	taskDistributor := NewRedisTaskDistributor(redisOpt)
	require.NotNil(t, taskDistributor)

	_, ok := taskDistributor.(*RedisTaskDistributor)
	require.True(t, ok)

	var _ TaskDistributor = taskDistributor
}
