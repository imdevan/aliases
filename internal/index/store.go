package index

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/aliases/internal/domain"
	"github.com/aliases/internal/utils"
)

// Store provides SQLite-backed alias persistence and search.
type Store struct {
	db     *sql.DB
	dbPath string
}

// NewStore opens (or creates) the SQLite index at the default XDG data path.
func NewStore() (*Store, error) {
	dbPath := DefaultDBPath()
	return NewStoreAt(dbPath)
}

// NewStoreAt opens (or creates) the SQLite index at a specific path.
func NewStoreAt(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("index: create data dir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("index: open db: %w", err)
	}

	// Enable WAL mode for concurrent reads during background refresh.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("index: set WAL: %w", err)
	}

	s := &Store{db: db, dbPath: dbPath}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// DefaultDBPath returns the standard location for the index database.
func DefaultDBPath() string {
	return filepath.Join(utils.XDGDataHome(), "aliases", "index.db")
}

// Close releases the database connection.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// migrate creates or updates the schema to the latest version.
func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS aliases (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT    NOT NULL,
			value       TEXT    NOT NULL,
			description TEXT    NOT NULL DEFAULT '',
			source      TEXT    NOT NULL DEFAULT '',
			mtime       INTEGER NOT NULL DEFAULT 0,
			global      INTEGER NOT NULL DEFAULT 0,
			suffix      INTEGER NOT NULL DEFAULT 0,
			indexed_at  TEXT    NOT NULL DEFAULT '',
			UNIQUE(name, source)
		);

		CREATE INDEX IF NOT EXISTS idx_aliases_name   ON aliases(name);
		CREATE INDEX IF NOT EXISTS idx_aliases_source ON aliases(source);

		CREATE TABLE IF NOT EXISTS source_files (
			path       TEXT    PRIMARY KEY,
			mtime      INTEGER NOT NULL DEFAULT 0,
			indexed_at TEXT    NOT NULL DEFAULT ''
		);
	`)
	if err != nil {
		return fmt.Errorf("index: migrate: %w", err)
	}
	return nil
}

// UpsertAlias inserts or updates a single alias in the index.
func (s *Store) UpsertAlias(a domain.Alias, mtime int64, global bool) error {
	suffix := strings.HasSuffix(a.Name, "-") || strings.HasSuffix(a.Name, "_")
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := s.db.Exec(`
		INSERT INTO aliases (name, value, description, source, mtime, global, suffix, indexed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(name, source) DO UPDATE SET
			value       = excluded.value,
			description = excluded.description,
			mtime       = excluded.mtime,
			global      = excluded.global,
			suffix      = excluded.suffix,
			indexed_at  = excluded.indexed_at
	`, a.Name, a.Value, a.Description, a.SourceFile, mtime, boolToInt(global), boolToInt(suffix), now)
	if err != nil {
		return fmt.Errorf("index: upsert alias %q: %w", a.Name, err)
	}
	return nil
}

// BulkUpsert efficiently inserts/updates a batch of aliases from one source.
func (s *Store) BulkUpsert(aliases []domain.Alias, source string, mtime int64, global bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("index: begin tx: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)
	stmt, err := tx.Prepare(`
		INSERT INTO aliases (name, value, description, source, mtime, global, suffix, indexed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(name, source) DO UPDATE SET
			value       = excluded.value,
			description = excluded.description,
			mtime       = excluded.mtime,
			global      = excluded.global,
			suffix      = excluded.suffix,
			indexed_at  = excluded.indexed_at
	`)
	if err != nil {
		return fmt.Errorf("index: prepare upsert: %w", err)
	}
	defer stmt.Close()

	for _, a := range aliases {
		a.SourceFile = source
		sfx := strings.HasSuffix(a.Name, "-") || strings.HasSuffix(a.Name, "_")
		if _, err := stmt.Exec(a.Name, a.Value, a.Description, source, mtime, boolToInt(global), boolToInt(sfx), now); err != nil {
			return fmt.Errorf("index: bulk upsert %q: %w", a.Name, err)
		}
	}

	// Remove aliases from this source that no longer exist in the file.
	nameSet := make(map[string]bool, len(aliases))
	for _, a := range aliases {
		nameSet[a.Name] = true
	}

	rows, err := tx.Query("SELECT id, name FROM aliases WHERE source = ?", source)
	if err != nil {
		return fmt.Errorf("index: query stale: %w", err)
	}
	var staleIDs []int64
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			rows.Close()
			return err
		}
		if !nameSet[name] {
			staleIDs = append(staleIDs, id)
		}
	}
	rows.Close()

	for _, id := range staleIDs {
		if _, err := tx.Exec("DELETE FROM aliases WHERE id = ?", id); err != nil {
			return fmt.Errorf("index: delete stale: %w", err)
		}
	}

	// Update source_files tracking.
	if _, err := tx.Exec(`
		INSERT INTO source_files (path, mtime, indexed_at)
		VALUES (?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET mtime = excluded.mtime, indexed_at = excluded.indexed_at
	`, source, mtime, now); err != nil {
		return fmt.Errorf("index: update source_files: %w", err)
	}

	return tx.Commit()
}

// Search finds aliases matching a query string (case-insensitive substring match
// across name, value, and description).
func (s *Store) Search(query string) ([]domain.Alias, error) {
	pattern := "%" + query + "%"
	rows, err := s.db.Query(`
		SELECT name, value, description, source, global, suffix
		FROM aliases
		WHERE name LIKE ? OR value LIKE ? OR description LIKE ?
		ORDER BY name ASC
	`, pattern, pattern, pattern)
	if err != nil {
		return nil, fmt.Errorf("index: search: %w", err)
	}
	defer rows.Close()

	return scanAliases(rows)
}

// All returns every alias in the index.
func (s *Store) All() ([]domain.Alias, error) {
	rows, err := s.db.Query(`
		SELECT name, value, description, source, global, suffix
		FROM aliases
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("index: all: %w", err)
	}
	defer rows.Close()

	return scanAliases(rows)
}

// BySource returns aliases from a specific source file.
func (s *Store) BySource(source string) ([]domain.Alias, error) {
	rows, err := s.db.Query(`
		SELECT name, value, description, source, global, suffix
		FROM aliases
		WHERE source = ?
		ORDER BY name ASC
	`, source)
	if err != nil {
		return nil, fmt.Errorf("index: by source: %w", err)
	}
	defer rows.Close()

	return scanAliases(rows)
}

// SourceMtime returns the stored mtime for a source file, or 0 if not indexed.
func (s *Store) SourceMtime(path string) (int64, error) {
	var mtime int64
	err := s.db.QueryRow("SELECT mtime FROM source_files WHERE path = ?", path).Scan(&mtime)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("index: source mtime: %w", err)
	}
	return mtime, nil
}

// DeleteBySource removes all aliases from a given source file.
func (s *Store) DeleteBySource(source string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM aliases WHERE source = ?", source); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM source_files WHERE path = ?", source); err != nil {
		return err
	}
	return tx.Commit()
}

// Stats returns basic index statistics.
func (s *Store) Stats() (totalAliases int, totalSources int, err error) {
	if err = s.db.QueryRow("SELECT COUNT(*) FROM aliases").Scan(&totalAliases); err != nil {
		return
	}
	err = s.db.QueryRow("SELECT COUNT(*) FROM source_files").Scan(&totalSources)
	return
}

func scanAliases(rows *sql.Rows) ([]domain.Alias, error) {
	var aliases []domain.Alias
	for rows.Next() {
		var a domain.Alias
		var global, suffix int
		if err := rows.Scan(&a.Name, &a.Value, &a.Description, &a.SourceFile, &global, &suffix); err != nil {
			return nil, fmt.Errorf("index: scan: %w", err)
		}
		a.Global = global == 1
		a.Suffix = suffix == 1
		aliases = append(aliases, a)
	}
	return aliases, rows.Err()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
