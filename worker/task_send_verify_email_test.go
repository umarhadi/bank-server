package worker

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/umarhadi/bank-server/db/mock"
	db "github.com/umarhadi/bank-server/db/sqlc"
	mockemail "github.com/umarhadi/bank-server/mail/mock"
	"go.uber.org/mock/gomock"
)

func TestDistributeTaskSendVerifyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	redisOpt := asynq.RedisClientOpt{
		Addr: "localhost:6379",
	}
	taskDistributor := NewRedisTaskDistributor(redisOpt)

	tests := []struct {
		name       string
		payload    *PayloadSendVerifyEmail
		buildStubs func()
		checkErr   func(err error)
	}{
		{
			name: "InvalidPayload",
			payload: &PayloadSendVerifyEmail{
				Username: string([]byte{0xff, 0xfe, 0xfd}), // Invalid UTF-8
			},
			buildStubs: func() {},
			checkErr: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to marshal task payload")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs()
			err := taskDistributor.DistributeTaskSendVerifyEmail(context.Background(), tc.payload)
			tc.checkErr(err)
		})
	}
}

func TestProcessTaskSendVerifyEmailErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	emailSender := mockemail.NewMockEmailSender(ctrl)
	taskProcessor := NewRedisTaskProcessor(asynq.RedisClientOpt{}, store, emailSender)

	user := db.User{
		Username: "test_user",
		FullName: "Test User",
		Email:    "test@example.com",
	}

	tests := []struct {
		name       string
		payload    *PayloadSendVerifyEmail
		buildStubs func()
		checkErr   func(error)
	}{
		{
			name: "CreateVerifyEmailError",
			payload: &PayloadSendVerifyEmail{
				Username: user.Username,
			},
			buildStubs: func() {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateVerifyEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.VerifyEmail{}, errors.New("db error"))
			},
			checkErr: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to create verify email")
			},
		},
		{
			name: "SendEmailError",
			payload: &PayloadSendVerifyEmail{
				Username: user.Username,
			},
			buildStubs: func() {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateVerifyEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.VerifyEmail{
						Username:   user.Username,
						Email:      user.Email,
						SecretCode: "secret-code",
					}, nil)

				emailSender.EXPECT().
					SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("failed to send email"))
			},
			checkErr: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to send verify email")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs()

			taskPayload, err := json.Marshal(tc.payload)
			require.NoError(t, err)

			task := asynq.NewTask(TaskSendVerifyEmail, taskPayload)
			err = taskProcessor.ProcessTaskSendVerifyEmail(context.Background(), task)

			tc.checkErr(err)
		})
	}
}
