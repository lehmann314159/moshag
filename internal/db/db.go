package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps a *sql.DB connection to the SQLite database.
type DB struct {
	*sql.DB
}

// Open opens (or creates) the SQLite database at path and runs schema migrations.
func Open(path string) (*DB, error) {
	// Ensure the parent directory exists.
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db directory: %w", err)
		}
	}

	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// Enable WAL mode for better concurrent read performance.
	if _, err := sqlDB.Exec("PRAGMA journal_mode=WAL"); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("set WAL mode: %w", err)
	}

	db := &DB{sqlDB}
	if err := db.migrate(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

// migrate runs CREATE TABLE IF NOT EXISTS statements.
func (db *DB) migrate() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			provider    TEXT NOT NULL,
			provider_id TEXT NOT NULL,
			name        TEXT NOT NULL,
			email       TEXT NOT NULL DEFAULT '',
			avatar_url  TEXT NOT NULL DEFAULT '',
			created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(provider, provider_id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS adventures (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id      INTEGER NOT NULL REFERENCES users(id),
			title        TEXT NOT NULL DEFAULT 'Untitled Adventure',
			mode         TEXT NOT NULL DEFAULT 'manual',
			current_step TEXT NOT NULL DEFAULT 'scenario',
			state_json   TEXT NOT NULL DEFAULT '{}',
			created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			adventure_id INTEGER NOT NULL REFERENCES adventures(id),
			role         TEXT NOT NULL,
			step         TEXT NOT NULL,
			content      TEXT NOT NULL,
			created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	return db.ensureGuestUser()
}

// GuestUserID is the DB user ID used for unauthenticated sessions.
const GuestUserID int64 = 1

// ensureGuestUser creates the guest user if it doesn't exist yet.
func (db *DB) ensureGuestUser() error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO users (id, provider, provider_id, name, email)
		VALUES (1, 'guest', 'guest', 'Guest', '')
	`)
	return err
}
