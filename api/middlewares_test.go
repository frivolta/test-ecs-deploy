package api

import (
	mockdb "birdie/db/mock"
	db "birdie/db/sqlc"
	"birdie/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func generateBearer(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}

func TestFirebaseAuth(t *testing.T) {
	testCase := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore, uuid string)
		setupAuth     func(request *http.Request, config util.Config)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Auth user already present in db",
			buildStubs: func(store *mockdb.MockStore, uuid string) {
				store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{
					ID:        1,
					FullName:  uuid,
					Role:      db.RoleTEACHER,
					Email:     "user1@test.com",
					UpdatedAt: time.Time{},
				}, nil)
			},
			setupAuth: func(request *http.Request, cfg util.Config) {
				token, err := util.GetIdTokenFromUUID(cfg.TestUUID, &cfg)
				require.NoError(t, err)
				bearer := util.GenerateBearerString(token)
				request.Header.Set("Authorization", bearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "User not in db has default role",
			buildStubs: func(store *mockdb.MockStore, uuid string) {
				store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{}, sql.ErrNoRows)
				store.EXPECT().CreateUser(gomock.Any(), db.CreateUserParams{
					FullName: uuid,
					Role:     db.RoleTEACHER,
					Email:    "user1@test.com",
				}).Times(1).Return(db.User{
					ID:        1,
					FullName:  uuid,
					Role:      db.RoleTEACHER,
					Email:     "user1@test.com",
					UpdatedAt: time.Time{},
				}, nil)
			},
			setupAuth: func(request *http.Request, cfg util.Config) {
				token, err := util.GetIdTokenFromUUID(cfg.TestUUID, &cfg)
				require.NoError(t, err)
				bearer := util.GenerateBearerString(token)
				request.Header.Set("Authorization", bearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Cannot create new user",
			buildStubs: func(store *mockdb.MockStore, uuid string) {
				store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{}, sql.ErrNoRows)
				store.EXPECT().CreateUser(gomock.Any(), db.CreateUserParams{
					FullName: uuid,
					Role:     db.RoleTEACHER,
					Email:    "user1@test.com",
				}).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			setupAuth: func(request *http.Request, cfg util.Config) {
				token, err := util.GetIdTokenFromUUID(cfg.TestUUID, &cfg)
				require.NoError(t, err)
				bearer := util.GenerateBearerString(token)
				request.Header.Set("Authorization", bearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				data, e := io.ReadAll(recorder.Body)
				require.NoError(t, e)
				var res map[string]string
				err := json.Unmarshal(data, &res)
				require.NoError(t, err)
				require.Equal(t, res, map[string]string{"error": "database connection error"})
			},
		},
		{
			name: "Not authorized user gets error 403",
			buildStubs: func(store *mockdb.MockStore, uuid string) {
				store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{
					ID:        1,
					FullName:  uuid,
					Role:      db.RolePARENT,
					Email:     "user1@test.com",
					UpdatedAt: time.Time{},
				}, nil)
			},
			setupAuth: func(request *http.Request, cfg util.Config) {
				token, err := util.GetIdTokenFromUUID(cfg.TestUUID, &cfg)
				require.NoError(t, err)
				bearer := util.GenerateBearerString(token)
				request.Header.Set("Authorization", bearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Returns 401 with invalid token",
			buildStubs: func(store *mockdb.MockStore, uuid string) {
			},
			setupAuth: func(request *http.Request, cfg util.Config) {
				token, err := util.GetIdTokenFromUUID(cfg.TestUUID, &cfg)
				require.NoError(t, err)
				bearer := util.GenerateBearerString(token + "invalid")
				request.Header.Set("Authorization", bearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				data, e := io.ReadAll(recorder.Body)
				require.NoError(t, e)
				var res map[string]string
				err := json.Unmarshal(data, &res)
				require.NoError(t, err)
				require.Equal(t, res, map[string]string{"error": "invalid token"})
			},
		},
		{
			name: "Returns 401 empty token",
			buildStubs: func(store *mockdb.MockStore, uuid string) {
			},
			setupAuth: func(request *http.Request, cfg util.Config) {
				_, err := util.GetIdTokenFromUUID(cfg.TestUUID, &cfg)
				require.NoError(t, err)
				bearer := util.GenerateBearerString("")
				request.Header.Set("Authorization", bearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				data, e := io.ReadAll(recorder.Body)
				require.NoError(t, e)
				var res map[string]string
				err := json.Unmarshal(data, &res)
				require.NoError(t, err)
				require.Equal(t, res, map[string]string{"error": "empty token"})
			},
		},
	}
	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/scoped")
			server, cfg := newTestServer(t, store)
			tc.buildStubs(store, cfg.TestUUID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			tc.setupAuth(request, cfg)
			require.NoError(t, err)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
