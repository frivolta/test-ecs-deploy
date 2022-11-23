package api

import (
	mockdb "birdie/db/mock"
	db "birdie/db/sqlc"
	"birdie/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateKidApi(t *testing.T) {
	kid := rndKid()
	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(request *http.Request, config util.Config)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":    kid.Name,
				"surname": kid.Surname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateKidParams{
					Name:    kid.Name,
					Surname: kid.Surname,
				}
				store.EXPECT().CreateKid(gomock.Any(), gomock.Eq(arg)).Times(1).Return(kid, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchKid(t, recorder.Body, kid)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			body: gin.H{
				"name":    kid.Name,
				"surname": kid.Surname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateKidParams{
					Name:    kid.Name,
					Surname: kid.Surname,
				}
				store.EXPECT().CreateKid(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Kid{}, sql.ErrConnDone)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Validation error",
			body: gin.H{
				"name":    "kid.Name",
				"surname": kid.Surname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateKid(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server, cfg := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := "/api/v1/kids/"
			store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{
				ID:        1,
				FullName:  cfg.TestUUID,
				Role:      db.RoleTEACHER,
				Email:     "user1@test.com",
				UpdatedAt: time.Time{},
			}, nil)
			bd, e := json.Marshal(tc.body)
			require.NoError(t, e)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bd))
			require.NoError(t, err)
			tc.setupAuth(request, cfg)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetKidAPI(t *testing.T) {
	kid := rndKid()
	testCase := []struct {
		name          string
		kidID         int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(request *http.Request, config util.Config)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "Ok",
			kidID: kid.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetKid(gomock.Any(), gomock.Eq(kid.ID)).Times(1).Return(kid, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchKid(t, recorder.Body, kid)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Not found",
			kidID: kid.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetKid(gomock.Any(), kid.ID).Times(1).Return(db.Kid{}, sql.ErrNoRows)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "Internal server error",
			kidID: kid.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetKid(gomock.Any(), gomock.Any()).Times(1).Return(db.Kid{}, sql.ErrConnDone)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "Invalid ID",
			kidID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetKid(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server, cfg := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/kids/%d", tc.kidID)
			store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{
				ID:        1,
				FullName:  cfg.TestUUID,
				Role:      db.RoleTEACHER,
				Email:     "user1@test.com",
				UpdatedAt: time.Time{},
			}, nil)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(request, cfg)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAllKidsAPI(t *testing.T) {
	n := 5
	kids := make([]db.Kid, n)
	for i := 0; i < n; i++ {
		kids[i] = rndKid()
	}
	testCase := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(request *http.Request, config util.Config)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllKids(gomock.Any()).Times(1).Return(kids, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllKids(gomock.Any()).Times(1).Return([]db.Kid{}, sql.ErrConnDone)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server, cfg := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/kids/")
			store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{
				ID:        1,
				FullName:  cfg.TestUUID,
				Role:      db.RoleTEACHER,
				Email:     "user1@test.com",
				UpdatedAt: time.Time{},
			}, nil)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(request, cfg)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func rndKid() db.Kid {
	return db.Kid{
		ID:      util.RandomInt(1, 1000),
		Name:    util.RandomName(),
		Surname: util.RandomName(),
	}
}
func requireBodyMatchKid(t *testing.T, body *bytes.Buffer, k db.Kid) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotKid db.Kid
	err = json.Unmarshal(data, &gotKid)
	require.NoError(t, err)
	require.Equal(t, k, gotKid)
}
