//go:build sqlite

package storage

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
)

const (
	sqlUPSERT = `INSERT INTO chat_user (chat_id, user_id) VALUES (?, ?)
    ON CONFLICT(chat_id, user_id) DO UPDATE 
	SET value = value + 1, updated_at = current_timestamp
	;`
	sqlSELECT = `SELECT value
	FROM chat_user
	WHERE chat_id = ? and user_id = ?
	;`
	sqlRESET = `UPDATE chat_user SET value = 0
	WHERE chat_id = ? and user_id = ?
	;`
)

//go:embed migrations/sqlite/*.sql
var embedMigrations embed.FS

func up(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set dialect: %v", err)
	}

	if err := goose.Up(db, "migrations/sqlite"); err != nil {
		return fmt.Errorf("up: %v", err)
	}

	return nil
}

func (s *anyStorage) Inc(chatID, userID int64) error {
	_, err := s.pool.Exec(sqlUPSERT, chatID, userID)
	if err != nil {
		return fmt.Errorf("upsert: %v", err)
	}

	return nil
}

func (s *anyStorage) Get(chatID, userID int64) (int64, error) {
	var value int64

	if err := s.pool.QueryRow(sqlSELECT, chatID, userID).Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("select: %v", err)
	}

	return value, nil
}

func (s *anyStorage) Reset(chatID, userID int64) error {
	_, err := s.pool.Exec(sqlRESET, chatID, userID)

	return err
}
