package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	if dsn == "" {
		return nil, errors.New("dsn required")
	}

	if dsn != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(dsn), 0700); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open database")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "cannot ping database")
	}

	if _, err = db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, errors.Wrap(err, "foreign keys pragma")
	}

	s := &Storage{
		db: db,
	}

	return s, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Migrate(dir string) error {
	if dir == "" {
		return errors.New("dir is required")
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(s.db, dir); err != nil {
		return err
	}

	return nil
}
