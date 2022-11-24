package api

import (
	db "birdie/db/sqlc"
	"birdie/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func newTestServer(t *testing.T, store db.Store) (*Server, util.Config) {
	realConfig, err := util.LoadConfig("..", "../serviceAccountKey.json")
	if err != nil {
		log.Fatal("Cannot load configuration", err)
	}

	server, err := NewServer(store, realConfig)
	require.NoError(t, err)

	return server, realConfig
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
