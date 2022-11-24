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

func TestUpdateTeacherNote(t *testing.T) {
	tc := rndTeacher()
	nt, _ := createTeacherNote(t, tc.ID)
	newNote := util.RandomString(5)
	upd := nt
	upd.Note = newNote
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
				"note": newNote,
			},
			query: nt.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateNoteParams{
					ID:   nt.ID,
					Note: newNote,
					Date: nt.Date,
				}
				store.EXPECT().GetTeacherNote(gomock.Any(), nt.ID).Times(1).Return(nt, nil)
				store.EXPECT().UpdateNote(gomock.Any(), arg).Times(1).Return(upd, nil)
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
				arg := db.UpdateNoteParams{
					ID:   nt.ID,
					Note: newNote,
					Date: nt.Date,
				}
				store.EXPECT().GetTeacherNote(gomock.Any(), nt.ID).Times(0)
				store.EXPECT().UpdateNote(gomock.Any(), arg).Times(0)
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
			url := fmt.Sprintf("/api/v1/teacher_notes/%d", tc.query)
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

func TestCreateTeacherNotes(t *testing.T) {
	tc := rndTeacher()
	nt, dt := createTeacherNote(t, tc.ID)
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
				"note": nt.Note,
				"date": dt,
			},
			query: tc.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTeacherNoteParams{
					Note: nt.Note,
					TeacherID: sql.NullInt64{
						Int64: tc.ID,
						Valid: true,
					},
					Date: nt.Date,
				}
				store.EXPECT().CreateTeacherNote(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nt, nil)
				store.EXPECT().GetTeacher(gomock.Any(), tc.ID).Times(1).Return(tc, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireTeacherNotesBodyMatcher(t, recorder.Body, nt)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid body",
			body: gin.H{
				"note": "",
				"date": "",
			},
			query: tc.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateTeacherNote(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetTeacher(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				fmt.Println(recorder.Body)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Teacher does not exists",
			body: gin.H{
				"note": nt.Note,
				"date": dt,
			},
			query: tc.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacher(gomock.Any(), tc.ID).Times(1).Return(db.Teacher{}, sql.ErrNoRows)
				store.EXPECT().CreateTeacherNote(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				fmt.Println(recorder.Body)
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			body: gin.H{
				"note": nt.Note,
				"date": dt,
			},
			query: tc.ID,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTeacherNoteParams{
					Note: nt.Note,
					TeacherID: sql.NullInt64{
						Int64: tc.ID,
						Valid: true,
					},
					Date: nt.Date,
				}
				store.EXPECT().CreateTeacherNote(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.TeacherNote{}, sql.ErrConnDone)
				store.EXPECT().GetTeacher(gomock.Any(), tc.ID).Times(1).Return(tc, nil)
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
			url := fmt.Sprintf("/api/v1/teacher_notes/%d", tc.query)
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

func TestGetTeacherNoteByDate(t *testing.T) {
	tc := rndTeacher()
	tn, dt := createTeacherNote(t, tc.ID)
	cnv, _ := util.ConvertDate(dt)
	res := []db.GetTeacherNotesByDateRow{
		{
			ID:   tn.ID,
			Note: tn.Note,
			Date: tn.Date,
			TeacherName: sql.NullString{
				String: tc.Name,
				Valid:  true,
			},
			TeacherSurname: sql.NullString{
				String: tc.Surname,
				Valid:  true,
			},
			TeacherID: sql.NullInt64{
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
				"date": dt,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacherNotesByDate(gomock.Any(), cnv).Times(1).Return(res, nil)
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
				"date": dt,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacherNotesByDate(gomock.Any(), gomock.Any()).Times(1).Return([]db.GetTeacherNotesByDateRow{}, sql.ErrConnDone)
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
				store.EXPECT().GetTeacherNotesByDate(gomock.Any(), gomock.Any()).Times(0)
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
			request, err := http.NewRequest(http.MethodPost, "/api/v1/teacher_notes/date", bytes.NewReader(bd))
			require.NoError(t, err)
			tc.setupAuth(request, cfg)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetTeacherNotesByPeriod(t *testing.T) {
	tc := rndTeacher()
	tn, dt := createTeacherNote(t, tc.ID)
	cnv, _ := util.ConvertDate(dt)
	res := []db.GetTeacherNotesByPeriodRow{
		{
			ID:   tn.ID,
			Note: tn.Note,
			Date: tn.Date,
			TeacherName: sql.NullString{
				String: tc.Name,
				Valid:  true,
			},
			TeacherSurname: sql.NullString{
				String: tc.Surname,
				Valid:  true,
			},
			TeacherID: sql.NullInt64{
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

				arg := db.GetTeacherNotesByPeriodParams{
					Date:   cnv,
					Date_2: cnv,
				}
				store.EXPECT().GetTeacherNotesByPeriod(gomock.Any(), arg).Times(1).Return(res, nil)
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
				store.EXPECT().GetTeacherNotesByPeriod(gomock.Any(), gomock.Any()).Times(1).Return([]db.GetTeacherNotesByPeriodRow{}, sql.ErrConnDone)
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
				store.EXPECT().GetTeacherNotesByPeriod(gomock.Any(), gomock.Any()).Times(0)
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
			request, err := http.NewRequest(http.MethodPost, "/api/v1/teacher_notes/period", bytes.NewReader(bd))
			require.NoError(t, err)
			tc.setupAuth(request, cfg)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAllTeacher(t *testing.T) {
	n := 5
	teachers := make([]db.Teacher, n)
	teacherNotes := make([]db.TeacherNote, n)
	for i := 0; i < n; i++ {
		teachers[i] = rndTeacher()
		teacherNotes[i], _ = createTeacherNote(t, teachers[i].ID)
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
				store.EXPECT().GetAllTeacherNotes(gomock.Any()).Times(1).Return(teacherNotes, nil)
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
				store.EXPECT().GetAllTeacherNotes(gomock.Any()).Times(1).Return([]db.TeacherNote{}, sql.ErrNoRows)
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
				store.EXPECT().GetAllTeacherNotes(gomock.Any()).Times(1).Return([]db.TeacherNote{}, sql.ErrConnDone)
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
			url := fmt.Sprintf("/api/v1/teacher_notes/")
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

// Create a mock note with converted date from a string
func createTeacherNote(t *testing.T, tID int64) (db.TeacherNote, string) {
	dateString := fmt.Sprintf("%d-%d-%d", 2022, util.RandomInt(10, 12), util.RandomInt(10, 28))
	convertedDate, err := util.ConvertDate(dateString)
	require.NoError(t, err)
	note := db.TeacherNote{
		ID:   util.RandomInt(1, 1000),
		Note: util.RandomString(150),
		TeacherID: sql.NullInt64{
			Int64: tID,
			Valid: true,
		},
		Date: convertedDate,
	}
	return note, dateString
}

// Require body to match multiple notes
func requireMultipleTeacherNotesBodyMatcher(t *testing.T, body *bytes.Buffer, nt []db.TeacherNote) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotTn db.TeacherNote
	err = json.Unmarshal(data, &gotTn)
	require.NoError(t, err)
	require.Equal(t, nt, gotTn)
}

// Require the body to match the json result
func requireTeacherNotesBodyMatcher(t *testing.T, body *bytes.Buffer, nt db.TeacherNote) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotTn db.TeacherNote
	err = json.Unmarshal(data, &gotTn)
	require.NoError(t, err)
	require.Equal(t, nt, gotTn)
}
