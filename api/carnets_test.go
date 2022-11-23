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

func rndCarnet(t *testing.T, k db.Kid) (db.Carnet, string) {
	d := "2022-10-15"
	cnv, e := util.ConvertDate(d)
	require.NoError(t, e)
	return db.Carnet{
		ID:       util.RandomInt(1, 1000),
		Date:     cnv,
		Quantity: util.RandomInt32(1, 5),
		KidID:    k.ID,
	}, d
}

func TestCreateCarnetApi(t *testing.T) {
	kid := rndKid()
	c, d := rndCarnet(t, kid)
	cnv, e := util.ConvertDate(d)
	require.NoError(t, e)
	testCase := []struct {
		name          string
		body          gin.H
		query         int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(request *http.Request, config util.Config)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: kid.ID,
			body: gin.H{
				"date":     d,
				"quantity": c.Quantity,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCarnetParams{
					Date:     cnv,
					Quantity: c.Quantity,
					KidID:    c.KidID,
				}
				store.EXPECT().GetKid(gomock.Any(), gomock.Eq(arg.KidID)).Times(1).Return(kid, nil)
				store.EXPECT().CreateCarnet(gomock.Any(), gomock.Eq(arg)).Times(1).Return(c, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchCarnet(t, recorder.Body, c)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Validation error",
			body: gin.H{
				"date":     "invalid",
				"quantity": "invalid",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetKid(gomock.Any(), gomock.Eq(kid.ID)).Times(0)
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
			url := fmt.Sprintf("/api/v1/carnets/%d", tc.query)
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

func TestGetAllCarnetsAPi(t *testing.T) {
	n := 5
	cs := make([]db.GetAllCarnetsRow, n)
	k := rndKid()
	for i := 0; i < n; i++ {
		c, d := rndCarnet(t, k)
		cnv, e := util.ConvertDate(d)
		require.NoError(t, e)
		cs[i] = db.GetAllCarnetsRow{
			ID:       c.ID,
			Date:     cnv,
			Quantity: c.Quantity,
			KidName: sql.NullString{
				String: k.Name,
				Valid:  true,
			},
			KidSurname: sql.NullString{
				String: k.Surname,
				Valid:  true,
			},
			KidID: sql.NullInt64{
				Int64: k.ID,
				Valid: true,
			},
		}
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
				store.EXPECT().GetAllKids(gomock.Any()).Times(1).Return([]db.Kid{{
					ID:      1,
					Name:    "Name",
					Surname: "Surname",
				}}, nil)
				store.EXPECT().GetAllKidNotes(gomock.Any()).Times(1).Return([]db.KidNote{{
					ID:       1,
					Note:     "note",
					KidID:    1,
					Presence: nil,
					HasMeal:  true,
					Date:     time.Time{},
				}}, nil)
				store.EXPECT().GetAllCarnets(gomock.Any()).Times(1).Return(cs, nil)
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
				store.EXPECT().GetAllCarnets(gomock.Any()).Times(1).Return([]db.GetAllCarnetsRow{}, sql.ErrConnDone)
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
			url := fmt.Sprintf("/api/v1/carnets/")
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

func TestUpdateCarnets(t *testing.T) {
	k := rndKid()
	c, d := rndCarnet(t, k)
	cnv, _ := util.ConvertDate(d)
	upd := c
	upd.Quantity = 1
	upd.Date = cnv
	testCase := []struct {
		name          string
		body          gin.H
		query         int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(request *http.Request, config util.Config)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"date":     d,
				"quantity": upd.Quantity,
			},
			query: c.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateCarnetParams{
					ID:       c.ID,
					Date:     upd.Date,
					Quantity: upd.Quantity,
				}
				store.EXPECT().GetCarnet(gomock.Any(), arg.ID).Times(1).Return(c, nil)
				store.EXPECT().UpdateCarnet(gomock.Any(), arg).Times(1).Return(upd, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Invalid params",
			body:  gin.H{},
			query: c.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacherNote(gomock.Any(), c.ID).Times(0)
				store.EXPECT().UpdateNote(gomock.Any(), gomock.Any()).Times(0)
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
			url := fmt.Sprintf("/api/v1/carnets/%d", tc.query)
			store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{
				ID:        1,
				FullName:  cfg.TestUUID,
				Role:      db.RoleTEACHER,
				Email:     "user1@test.com",
				UpdatedAt: time.Time{},
			}, nil)
			bd, e := json.Marshal(tc.body)
			require.NoError(t, e)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(bd))
			require.NoError(t, err)
			tc.setupAuth(request, cfg)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchCarnet(t *testing.T, body *bytes.Buffer, c db.Carnet) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotCarnet db.Carnet
	err = json.Unmarshal(data, &gotCarnet)
	require.NoError(t, err)
	require.Equal(t, c, gotCarnet)
}
