package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	migration "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	uri string
	db *sql.DB
}

func DbConnect(uri string) (*Storage, error) {
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return &Storage{
		uri: uri,
		db: db,
	}, nil
}

func (s *Storage) Migrate() error {
	cfg := &migration.Config{
	}
	db, err := migration.WithInstance(s.db, cfg)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "", db)
	if err != nil {
		return err
	}
	return m.Up()
}
