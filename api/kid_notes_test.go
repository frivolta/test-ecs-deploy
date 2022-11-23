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

func TestUpdateKidNote(t *testing.T) {
	tc := rndKid()
	nt, _ := createKidNote(t, tc.ID)
	newNote := util.RandomString(5)
	upd := nt
	upd.Note = newNote
	upd.Presence = []db.Presence{"ABSENT"}
	upd.HasMeal = true
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
				"note":     newNote,
				"presence": upd.Presence,
				"has_meal": upd.HasMeal,
			},
			query: nt.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateKidNoteParams{
					ID:       nt.ID,
					Note:     newNote,
					Presence: upd.Presence,
					HasMeal:  upd.HasMeal,
				}
				store.EXPECT().GetKidNote(gomock.Any(), nt.ID).Times(1).Return(nt, nil)
				store.EXPECT().UpdateKidNote(gomock.Any(), arg).Times(1).Return(upd, nil)
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
			query: nt.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateKidNoteParams{
					ID:       nt.ID,
					Note:     newNote,
					Presence: upd.Presence,
					HasMeal:  upd.HasMeal,
				}
				store.EXPECT().GetKidNote(gomock.Any(), nt.ID).Times(0)
				store.EXPECT().UpdateKidNote(gomock.Any(), arg).Times(0)
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
			url := fmt.Sprintf("/api/v1/kid_notes/%d", tc.query)
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

// Create a mock note with converted date from a string
func createKidNote(t *testing.T, kID int64) (db.KidNote, string) {
	dateString := fmt.Sprintf("%d-%d-%d", 2022, util.RandomInt(10, 12), util.RandomInt(10, 28))
	convertedDate, err := util.ConvertDate(dateString)
	require.NoError(t, err)
	note := db.KidNote{
		ID:       util.RandomInt(1, 1000),
		Note:     util.RandomString(150),
		KidID:    kID,
		Date:     convertedDate,
		HasMeal:  false,
		Presence: []db.Presence{"MORNING", "AFTERNOON"},
	}
	return note, dateString
}

func TestCreateKidNotes(t *testing.T) {
	k := rndKid()
	nt, dt := createKidNote(t, k.ID)
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
				"note":     nt.Note,
				"date":     dt,
				"presence": nt.Presence,
				"has_meal": nt.HasMeal,
			},
			query: k.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateKidNoteParams{
					Note:     nt.Note,
					KidID:    k.ID,
					Date:     nt.Date,
					Presence: nt.Presence,
					HasMeal:  nt.HasMeal,
				}
				store.EXPECT().GetKidNotesByPeriod(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)
				store.EXPECT().GetKid(gomock.Any(), k.ID).Times(1).Return(k, nil)
				store.EXPECT().CreateKidNote(gomock.Any(), arg).Times(1).Return(nt, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireKidNotesBodyMatcher(t, recorder.Body, nt)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Kid already exits error",
			body: gin.H{
				"note":     nt.Note,
				"date":     dt,
				"presence": nt.Presence,
				"has_meal": nt.HasMeal,
			},
			query: k.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetKidNotesByPeriod(gomock.Any(), gomock.Any()).Times(1).Return([]db.GetKidNotesByPeriodRow{
					{
						ID:         0,
						Note:       "",
						Date:       time.Time{},
						Presence:   nil,
						HasMeal:    false,
						KidName:    sql.NullString{},
						KidSurname: sql.NullString{},
						KidID: sql.NullInt64{
							Int64: k.ID,
							Valid: true,
						},
					},
				}, nil)
				store.EXPECT().GetKid(gomock.Any(), k.ID).Times(1).Return(k, nil)
				store.EXPECT().CreateKidNote(gomock.Any(), gomock.Any()).Times(0)
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
			url := fmt.Sprintf("/api/v1/kid_notes/%d", tc.query)
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

func TestGetKidNotesByPeriodAPI(t *testing.T) {
	tc := rndKid()
	tn, dt := createKidNote(t, tc.ID)
	cnv, _ := util.ConvertDate(dt)
	res := []db.GetKidNotesByPeriodRow{
		{
			ID:       tc.ID,
			Note:     tn.Note,
			Date:     tn.Date,
			Presence: []db.Presence{"ABSENT"},
			HasMeal:  true,
			KidName: sql.NullString{
				String: tc.Name,
				Valid:  true,
			},
			KidSurname: sql.NullString{
				String: tc.Surname,
				Valid:  true,
			},
			KidID: sql.NullInt64{
				Int64: tc.ID,
				Valid: true,
			},
		},
	}

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
				"date1": dt,
				"date2": dt,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.GetKidNotesByPeriodParams{
					Date:   cnv,
					Date_2: cnv,
				}
				store.EXPECT().GetKidNotesByPeriod(gomock.Any(), arg).Times(1).Return(res, nil)
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
			body: gin.H{
				"date1": dt,
				"date2": dt,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetKidNotesByPeriod(gomock.Any(), gomock.Any()).Times(1).Return([]db.GetKidNotesByPeriodRow{}, sql.ErrConnDone)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetKidNotesByPeriod(gomock.Any(), gomock.Any()).Times(0)
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
			store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{
				ID:        1,
				FullName:  cfg.TestUUID,
				Role:      db.RoleTEACHER,
				Email:     "user1@test.com",
				UpdatedAt: time.Time{},
			}, nil)
			bd, e := json.Marshal(tc.body)
			require.NoError(t, e)
			request, err := http.NewRequest(http.MethodPost, "/api/v1/kid_notes/period", bytes.NewReader(bd))
			require.NoError(t, err)
			tc.setupAuth(request, cfg)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAllKidNotesAPi(t *testing.T) {
	n := 5
	teachers := make([]db.Kid, n)
	teacherNotes := make([]db.KidNote, n)
	for i := 0; i < n; i++ {
		teachers[i] = rndKid()
		teacherNotes[i], _ = createKidNote(t, teachers[i].ID)
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
				store.EXPECT().GetAllKidNotes(gomock.Any()).Times(1).Return(teacherNotes, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Not found",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllKidNotes(gomock.Any()).Times(1).Return([]db.KidNote{}, sql.ErrNoRows)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllKidNotes(gomock.Any()).Times(1).Return([]db.KidNote{}, sql.ErrConnDone)
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
			url := fmt.Sprintf("/api/v1/kid_notes/")
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

// Require the body to match the json result
func requireKidNotesBodyMatcher(t *testing.T, body *bytes.Buffer, nt db.KidNote) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotTn db.KidNote
	err = json.Unmarshal(data, &gotTn)
	require.NoError(t, err)
	require.Equal(t, nt, gotTn)
}
