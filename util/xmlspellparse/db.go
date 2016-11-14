package xmlspellparse

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// InitDB returns a *sqlx.DB initialized with connection parameters
// from this package's config.json file
// REMEMBER TO CLOSE THE CONNECTION WHEN DONE
func InitDB() (*sqlx.DB, error) {

	// Read config file into memory
	cFile, err := os.Open("config.json")
	cBytes, err := ioutil.ReadAll(cFile)
	cFile.Close()

	// Parse config file
	dbConfig := mysql.Config{}
	json.Unmarshal(cBytes, &dbConfig)
	if err != nil {
		return &sqlx.DB{}, err
	}

	db, err := sqlx.Connect("mysql", dbConfig.FormatDSN())
	if err != nil {
		return &sqlx.DB{}, err
	}
	return db, nil
}
