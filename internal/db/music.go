package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/KoLLlaka/poma-botv2.0/internal/model"
)

const (
	tableMusic = "musics"
)

type MusicStore interface {
	GetDuration(music *model.Playlist) error
	StoreDuration(music *model.Playlist) error
}

type musicStore struct {
	db *sql.DB
}

func NewMusicStore(db *sql.DB) MusicStore {
	return &musicStore{db: db}
}

func (m *musicStore) GetDuration(music *model.Playlist) error {
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE name = ?", tableMusic)
	if err := m.db.QueryRow(stmt, filepath.Base(music.Title)).Scan(&music.Title, &music.Duration); err != nil {
		return err
	}

	return nil
}

func (m *musicStore) StoreDuration(music *model.Playlist) error {
	stmt := fmt.Sprintf("INSERT INTO %s (name, duration) VALUES(?, ?)", tableMusic)
	if _, err := m.db.Exec(
		stmt,
		music.Title,
		music.Duration,
	); err != nil {
		return err
	}

	return nil
}
