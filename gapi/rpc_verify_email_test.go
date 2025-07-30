package gapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	mockdb "github.com/umarhadi/bank-server/db/mock"
	db "github.com/umarhadi/bank-server/db/sqlc"
	"github.com/umarhadi/bank-server/pb"
	"github.com/umarhadi/bank-server/util"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestVerifyEmailAPI(t *testing.T) {
	user, _ := randomUser(t, util.DepositorRole)
	validSecretCode := util.RandomString(32)
	
	testCases := []struct {
		name          string
		req           *pb.VerifyEmailRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.VerifyEmailResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.VerifyEmailRequest{
				EmailId:    1,
				SecretCode: validSecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.VerifyEmailTxResult{
						User: db.User{
							Username:        user.Username,
							IsEmailVerified: true,
						},
					}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.True(t, res.IsVerified)
			},
		},
		{
			name: "InternalError",
			req: &pb.VerifyEmailRequest{
				EmailId:    1,
				SecretCode: validSecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.VerifyEmailTxResult{}, fmt.Errorf("internal error"))
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "InvalidEmailId",
			req: &pb.VerifyEmailRequest{
				EmailId:    0,
				SecretCode: validSecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InvalidSecretCode",
			req: &pb.VerifyEmailRequest{
				EmailId:    1,
				SecretCode: "short",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil)
			res, err := server.VerifyEmail(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}