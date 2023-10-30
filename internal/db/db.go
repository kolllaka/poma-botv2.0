package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./"+dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err := migrationUp("file://static/db/migrations/", dsn); err != nil {
		return nil, err
	}

	return db, nil
}

func DownDB(dsn string) error {
	db, err := sql.Open("sqlite3", "./"+dsn)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	if err := migrationDown("file://static/db/migrations/", dsn); err != nil {
		return err
	}

	return nil
}

func migrationUp(migrationPath string, dsn string) error {
	migration, err := migrate.New(migrationPath, "sqlite3://"+dsn)
	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func migrationDown(migrationPath string, dsn string) error {
	migration, err := migrate.New(migrationPath, "sqlite3://"+dsn)
	if err != nil {
		return err
	}

	if err := migration.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func CloseDB(db *sql.DB) {
	db.Close()
}
