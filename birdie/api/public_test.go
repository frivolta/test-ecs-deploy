package api

import (
	mockdb "birdie/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	server, _ := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := "/"

	// Make request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	// Route
	server.Router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
}
