package api

import (
	mockdb "birdie/db/mock"
	db "birdie/db/sqlc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPrivateApi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	server, cfg := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := "/private"

	// Make request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	store.EXPECT().GetUser(gomock.Any(), "user1@test.com").Times(1).Return(db.User{
		ID:        1,
		FullName:  cfg.TestUUID,
		Role:      db.RoleTEACHER,
		Email:     "user1@test.com",
		UpdatedAt: time.Time{},
	}, nil)
	addAuthorization(request, cfg)

	// Route
	server.Router.ServeHTTP(recorder, request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, recorder.Code)
}
