package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/umarhadi/bank-server/db/mock"
	db "github.com/umarhadi/bank-server/db/sqlc"
	"github.com/umarhadi/bank-server/token"
)

type mockTokenMaker struct {
	token.Maker
	createTokenErr error
	callCount      int
	errOnCall      int
}

func (m *mockTokenMaker) CreateToken(username string, role string, duration time.Duration) (string, *token.Payload, error) {
	m.callCount++
	if m.errOnCall > 0 && m.callCount == m.errOnCall {
		return "", nil, m.createTokenErr
	}
	return m.Maker.CreateToken(username, role, duration)
}

func TestRenewAccessToken(t *testing.T) {
	user, _ := randomUser(t)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload)
		buildStubs    func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				refreshToken, refreshPayload, err := tokenMaker.CreateToken(
					user.Username,
					user.Role,
					time.Minute,
				)
				require.NoError(t, err)
				return refreshToken, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.ID)).
					Times(1).
					Return(db.Session{
						ID:           refreshPayload.ID,
						Username:     user.Username,
						RefreshToken: refreshToken,
						UserAgent:    "user-agent",
						ClientIp:     "127.0.0.1",
						IsBlocked:    false,
						ExpiresAt:    refreshPayload.ExpiredAt,
						CreatedAt:    time.Now(),
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var resp renewAccessTokenResponse
				require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))
				require.NotEmpty(t, resp.AccessToken)
				require.NotZero(t, resp.AccessTokenExpiresAt)
			},
		},
		{
			name: "NoRefreshToken",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				return "", nil
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ExpiredSession",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				refreshToken, refreshPayload, err := tokenMaker.CreateToken(
					user.Username,
					user.Role,
					time.Minute,
				)
				require.NoError(t, err)
				return refreshToken, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.ID)).
					Times(1).
					Return(db.Session{
						ID:           refreshPayload.ID,
						Username:     user.Username,
						RefreshToken: refreshToken,
						UserAgent:    "user-agent",
						ClientIp:     "127.0.0.1",
						IsBlocked:    false,
						ExpiresAt:    time.Now().Add(-time.Hour),
						CreatedAt:    time.Now(),
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "BlockedSession",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				refreshToken, refreshPayload, err := tokenMaker.CreateToken(
					user.Username,
					user.Role,
					time.Minute,
				)
				require.NoError(t, err)
				return refreshToken, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.ID)).
					Times(1).
					Return(db.Session{
						ID:           refreshPayload.ID,
						Username:     user.Username,
						RefreshToken: refreshToken,
						UserAgent:    "user-agent",
						ClientIp:     "127.0.0.1",
						IsBlocked:    true,
						ExpiresAt:    refreshPayload.ExpiredAt,
						CreatedAt:    time.Now(),
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidRefreshToken",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				return "invalid_token", nil
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SessionNotFound",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				refreshToken, refreshPayload, err := tokenMaker.CreateToken(
					user.Username,
					user.Role,
					time.Minute,
				)
				require.NoError(t, err)
				return refreshToken, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.ID)).
					Times(1).
					Return(db.Session{}, db.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				refreshToken, refreshPayload, err := tokenMaker.CreateToken(
					user.Username,
					user.Role,
					time.Minute,
				)
				require.NoError(t, err)
				return refreshToken, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.ID)).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "MismatchedUsername",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				refreshToken, refreshPayload, err := tokenMaker.CreateToken(
					user.Username,
					user.Role,
					time.Minute,
				)
				require.NoError(t, err)
				return refreshToken, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.ID)).
					Times(1).
					Return(db.Session{
						ID:           refreshPayload.ID,
						Username:     "different_user",
						RefreshToken: refreshToken,
						ExpiresAt:    refreshPayload.ExpiredAt,
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "MismatchedRefreshToken",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				refreshToken, refreshPayload, err := tokenMaker.CreateToken(
					user.Username,
					user.Role,
					time.Minute,
				)
				require.NoError(t, err)
				return refreshToken, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.ID)).
					Times(1).
					Return(db.Session{
						ID:           refreshPayload.ID,
						Username:     user.Username,
						RefreshToken: "different_token",
						ExpiresAt:    refreshPayload.ExpiredAt,
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "CreateTokenError",
			setupAuth: func(t *testing.T, tokenMaker token.Maker) (string, *token.Payload) {
				refreshToken, refreshPayload, err := tokenMaker.CreateToken(
					user.Username,
					user.Role,
					time.Minute,
				)
				require.NoError(t, err)
				return refreshToken, refreshPayload
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.ID)).
					Times(1).
					Return(db.Session{
						ID:           refreshPayload.ID,
						Username:     user.Username,
						RefreshToken: refreshToken,
						UserAgent:    "user-agent",
						ClientIp:     "127.0.0.1",
						IsBlocked:    false,
						ExpiresAt:    refreshPayload.ExpiredAt,
						CreatedAt:    time.Now(),
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := newTestServer(t, store)

			refreshToken, refreshPayload := tc.setupAuth(t, server.tokenMaker)

			if tc.name == "CreateTokenError" {
				server.tokenMaker = &mockTokenMaker{
					Maker:          server.tokenMaker,
					createTokenErr: errors.New("token creation error"),
					errOnCall:      1,
				}
			}

			tc.buildStubs(store, refreshToken, refreshPayload)

			url := "/tokens/renew_access"
			data, err := json.Marshal(gin.H{
				"refresh_token": refreshToken,
			})
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
