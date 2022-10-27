//go:build sqlite

package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"meBot/pkg/log"
)

func initDB(dbURL string) (*sql.DB, error) {
	if dbURL == "" {
		return nil, errors.New("empty db name")
	}

	log.Debug("connect to db...", zap.String("db_url", dbURL))

	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %v", err)
	}

	log.Debug("connection established")

	return db, nil
}
