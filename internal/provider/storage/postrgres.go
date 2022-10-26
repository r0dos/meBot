//go:build postgres

package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose"
)

//go:embed migrations/postgres/*.sql
var embedMigrations embed.FS

func up(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %v", err)
	}

	if err := goose.Up(db, "migrations/postgres"); err != nil {
		return fmt.Errorf("up: %v", err)
	}

	return nil
}
