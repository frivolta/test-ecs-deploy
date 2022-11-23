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

func TestCreateTeacherAPI(t *testing.T) {
	teacher := rndTeacher()
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
				"name":    teacher.Name,
				"surname": teacher.Surname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTeacherParams{
					Name:    teacher.Name,
					Surname: teacher.Surname,
				}
				store.EXPECT().CreateTeacher(gomock.Any(), gomock.Eq(arg)).Times(1).Return(teacher, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchTeacher(t, recorder.Body, teacher)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			body: gin.H{
				"name":    teacher.Name,
				"surname": teacher.Surname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTeacherParams{
					Name:    teacher.Name,
					Surname: teacher.Surname,
				}
				store.EXPECT().CreateTeacher(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Teacher{}, sql.ErrConnDone)
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
				"name":    "teacher.Name",
				"surname": teacher.Surname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateTeacher(gomock.Any(), gomock.Any()).Times(0)
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
			url := "/api/v1/teachers/"
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

func TestGetTeacherAPI(t *testing.T) {
	teacher := rndTeacher()
	testCase := []struct {
		name          string
		teacherID     int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(request *http.Request, config util.Config)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Ok",
			teacherID: teacher.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacher(gomock.Any(), gomock.Eq(teacher.ID)).Times(1).Return(teacher, nil)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchTeacher(t, recorder.Body, teacher)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Not found",
			teacherID: teacher.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacher(gomock.Any(), teacher.ID).Times(1).Return(db.Teacher{}, sql.ErrNoRows)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal server error",
			teacherID: teacher.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacher(gomock.Any(), gomock.Any()).Times(1).Return(db.Teacher{}, sql.ErrConnDone)
			},
			setupAuth: func(request *http.Request, config util.Config) {
				addAuthorization(request, config)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Invalid ID",
			teacherID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacher(gomock.Any(), gomock.Any()).Times(0)
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
			url := fmt.Sprintf("/api/v1/teachers/%d", tc.teacherID)
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

func TestGetAllTeacherAPI(t *testing.T) {
	n := 5
	teachers := make([]db.Teacher, n)
	for i := 0; i < n; i++ {
		teachers[i] = rndTeacher()
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
				store.EXPECT().GetAllTeachers(gomock.Any()).Times(1).Return(teachers, nil)
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
				store.EXPECT().GetAllTeachers(gomock.Any()).Times(1).Return([]db.Teacher{}, sql.ErrConnDone)
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
			url := fmt.Sprintf("/api/v1/teachers/")
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

func addAuthorization(request *http.Request, cfg util.Config) {
	token, _ := util.GetIdTokenFromUUID(cfg.TestUUID, &cfg)
	bearer := util.GenerateBearerString(token)
	request.Header.Set("Authorization", bearer)
}

func rndTeacher() db.Teacher {
	return db.Teacher{
		ID:      util.RandomInt(1, 1000),
		Name:    util.RandomName(),
		Surname: util.RandomName(),
	}
}
func requireBodyMatchTeacher(t *testing.T, body *bytes.Buffer, tc db.Teacher) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotTc db.Teacher
	err = json.Unmarshal(data, &gotTc)
	require.NoError(t, err)
	require.Equal(t, tc, gotTc)
}
