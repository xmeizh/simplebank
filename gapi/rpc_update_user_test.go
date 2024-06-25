package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/xmeizh/simplebank/db/mock"
	db "github.com/xmeizh/simplebank/db/postgresql"
	"github.com/xmeizh/simplebank/pb"
	"github.com/xmeizh/simplebank/token"
	"github.com/xmeizh/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(20)

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, resp *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
				Password: &newPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newFullName,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
				}

				expectedUser := db.User{
					Username:        user.Username,
					FullName:        newFullName,
					Email:           newEmail,
					CreatedAt:       user.CreatedAt,
					IsEmailVerified: user.IsEmailVerified,
				}

				store.EXPECT().UpdateUser(gomock.Any(), EqUpdateUserParams(arg, newPassword)).Times(1).Return(expectedUser, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				token, _, err := tokenMaker.CreateToken(user.Username, time.Minute)
				require.NoError(t, err)
				bearerToken := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				md := metadata.MD{
					authorizationHeader: []string{
						bearerToken,
					},
				}

				return metadata.NewIncomingContext(context.Background(), md)
			},
			checkResponse: func(t *testing.T, resp *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				createdUser := resp.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, newFullName, createdUser.FullName)
				require.Equal(t, newEmail, createdUser.Email)
			},
		},
		{
			name: "UserNotFound",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullName,
				Email:    &newEmail,
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), codes.NotFound)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)
			server := newTestServer(t, store, nil)

			ctx := tc.buildContext(t, server.tokenMaker)
			resp, err := server.UpdateUser(ctx, tc.req)

			tc.checkResponse(t, resp, err)
		})
	}
}

type eqUpdateUserParamsMatcher struct {
	arg      db.UpdateUserParams
	password string
}

func (expected eqUpdateUserParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	if actualArg.HashedPassword.Valid {
		err := util.CheckPassword(expected.password, actualArg.HashedPassword.String)

		if err != nil {
			return false
		}
		expected.arg.HashedPassword = actualArg.HashedPassword
		expected.arg.PasswordChangedAt = actualArg.PasswordChangedAt
	}

	return reflect.DeepEqual(expected.arg, actualArg)
}

func (expected eqUpdateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", expected.arg, expected.password)
}

func EqUpdateUserParams(arg db.UpdateUserParams, password string) gomock.Matcher {
	return eqUpdateUserParamsMatcher{arg, password}
}
