package db

import (
	"database/sql"
	"fmt"
	"go_final_project/pkg/constants"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("error opening database: " + err.Error())
	}

	if _, err := db.Exec(constants.Schema); err != nil {
		return fmt.Errorf("schema creating error: " + err.Error())
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
