//go:build postgres

package main

import (
	"database/sql"
	"errors"
)

func initDB(dbURL string) (*sql.DB, error) {
	return nil, errors.New("dont initialize")
}
