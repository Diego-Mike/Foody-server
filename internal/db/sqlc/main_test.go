package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Foody-App-Tech/Main-server/config"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := config.LoadEnv("../../../")
	if err != nil {
		log.Fatal("Can't load env variables :", err)
	}

	testDB, err = sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("Can't connect to db :", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
