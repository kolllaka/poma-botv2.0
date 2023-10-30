package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS musics (name TEXT UNIQUE, duration INTEGER NOT NULL);"); err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB(db *sql.DB) {
	db.Close()
}
