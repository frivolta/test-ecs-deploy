package db

import (
	"birdie/util"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..", "../../serviceAccountKey.json")

	if err != nil {
		log.Fatal("cannot load config", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the db", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
