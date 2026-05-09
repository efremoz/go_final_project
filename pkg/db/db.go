package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"go_final_project/pkg/constants"
)

var db *sql.DB

func Init(dbFile string) error {
	dir := filepath.Dir(dbFile)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("directory creation error %s: %w", dir, err)
		}
	}

	_, err := os.Stat(dbFile)

	db, err = sql.Open("sqlite", dbFile+"?cache=shared")
	if err != nil {
		return fmt.Errorf("error opening database: " + err.Error())
	}

	if _, err := db.Exec(constants.Schema); err != nil {
		_ = db.Close()
		return fmt.Errorf("schema creating error: " + err.Error())
	}

	db.SetMaxOpenConns(1)
	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
