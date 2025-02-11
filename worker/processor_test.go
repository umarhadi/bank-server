package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/umarhadi/bank-server/db/mock"
	db "github.com/umarhadi/bank-server/db/sqlc"
	mockemail "github.com/umarhadi/bank-server/mail/mock"
	"go.uber.org/mock/gomock"
)

func TestProcessTaskSendVerifyEmail(t *testing.T) {
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

	verifyEmail := db.VerifyEmail{
		ID:         1,
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: "secret-code",
	}

	// Expected email content
	subject := "Welcome to Very Bank Service"
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, verifyUrl)

	tests := []struct {
		name       string
		payload    *PayloadSendVerifyEmail
		buildStubs func()
		checkResp  func(err error)
	}{
		{
			name: "OK",
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
					Return(verifyEmail, nil)

				emailSender.EXPECT().
					SendEmail(
						subject,
						content,
						[]string{user.Email},
						[]string(nil),
						[]string(nil),
						[]string(nil),
					).
					Times(1).
					Return(nil)
			},
			checkResp: func(err error) {
				require.NoError(t, err)
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

			tc.checkResp(err)
		})
	}
}

func TestProcessTaskSendVerifyEmailError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	emailSender := mockemail.NewMockEmailSender(ctrl)
	taskProcessor := NewRedisTaskProcessor(asynq.RedisClientOpt{}, store, emailSender)

	user := db.User{
		Username: "test_user",
		Email:    "test@example.com",
	}

	tests := []struct {
		name       string
		payload    *PayloadSendVerifyEmail
		buildStubs func()
		checkResp  func(err error)
	}{
		{
			name: "GetUserError",
			payload: &PayloadSendVerifyEmail{
				Username: user.Username,
			},
			buildStubs: func() {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, errors.New("user not found"))
			},
			checkResp: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to get user")
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

			tc.checkResp(err)
		})
	}
}

func TestProcessorStartShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	// Allow any call to GetUser and return a dummy user.
	store.EXPECT().GetUser(gomock.Any(), gomock.Any()).AnyTimes().Return(
		db.User{
			Username: "dummy",
			FullName: "Dummy User",
			Email:    "dummy@example.com",
		}, nil,
	)
	// Allow any call to CreateVerifyEmail and return a dummy verify email.
	store.EXPECT().CreateVerifyEmail(gomock.Any(), gomock.Any()).AnyTimes().Return(
		db.VerifyEmail{
			ID:         1,
			Username:   "dummy",
			Email:      "dummy@example.com",
			SecretCode: "dummy-secret",
		}, nil,
	)

	// Optionally allow SendEmail calls.
	emailSender := mockemail.NewMockEmailSender(ctrl)
	emailSender.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().Return(nil)

	taskProcessor := NewRedisTaskProcessor(asynq.RedisClientOpt{}, store, emailSender)

	// Start the processor in a goroutine since it blocks.
	go func() {
		err := taskProcessor.Start()
		require.NoError(t, err)
	}()

	// Let the server start then immediately shutdown.
	time.Sleep(100 * time.Millisecond)
	taskProcessor.Shutdown()
}

func TestRedisTaskProcessorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	emailSender := mockemail.NewMockEmailSender(ctrl)

	// Create a test task that will fail
	taskPayload := []byte("invalid json")
	task := asynq.NewTask(TaskSendVerifyEmail, taskPayload)

	// Create processor with error handler
	processor := NewRedisTaskProcessor(asynq.RedisClientOpt{}, store, emailSender)

	// Process task - this should trigger error handler
	err := processor.ProcessTaskSendVerifyEmail(context.Background(), task)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal payload")
}
