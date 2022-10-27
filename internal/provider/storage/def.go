//go:build !sqlite && !postgres

package storage

import "database/sql"

func up(db *sql.DB) error {
	return nil
}

func inc(pool *sql.DB, chatID, userID int64) error {
	return nil
}

func get(pool *sql.DB, chatID, userID int64) (int64, error) {
	return 0, nil
}

func reset(pool *sql.DB, chatID, userID int64) error {
	return nil
}
