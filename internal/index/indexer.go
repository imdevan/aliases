package index

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliases/internal/alias"
	"github.com/aliases/internal/domain"
)

// Indexer scans alias source files and populates the SQLite index.
type Indexer struct {
	store         *Store
	aliasFile     string
	indexFolders  []string
	shell         string
	cacheInterval int // seconds
}

// NewIndexer creates an indexer that reads from the given alias file and
// index_folders glob patterns, writing results to the provided store.
func NewIndexer(store *Store, cfg domain.Config) *Indexer {
	return &Indexer{
		store:         store,
		aliasFile:     cfg.ResolvedAliasFile(),
		indexFolders:  cfg.IndexFolders,
		shell:         cfg.Shell,
		cacheInterval: cfg.CacheInterval,
	}
}

// IndexResult summarises one indexing run.
type IndexResult struct {
	FilesScanned  int
	FilesSkipped  int
	AliasesStored int
	Duration      time.Duration
	Errors        []error
}

// RefreshAsync performs background refresh if NeedsRefresh returns true.
// If detached is true, it spawns a background subprocess so shell startup is not blocked.
func (ix *Indexer) RefreshAsync(detached bool) error {
	if !ix.NeedsRefresh() {
		return nil
	}
	if detached {
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		cmd := exec.Command(exe, "index", "--bg")
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil
		return cmd.Start()
	}

	go ix.Refresh()
	return nil
}

// Refresh indexes all sources, skipping files whose mtime has not changed.
func (ix *Indexer) Refresh() IndexResult {
	start := time.Now()
	result := IndexResult{}

	// 1. Index the primary alias file.
	if ix.aliasFile != "" {
		if err := ix.indexFile(ix.aliasFile, true, &result); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("primary %s: %w", ix.aliasFile, err))
		}
	}

	// 2. Separate include patterns and "!" exclude patterns.
	var includes []string
	var excludes []string
	for _, raw := range ix.indexFolders {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "!") {
			excludes = append(excludes, strings.TrimPrefix(trimmed, "!"))
		} else {
			includes = append(includes, trimmed)
		}
	}

	visited := make(map[string]bool)
	if ix.aliasFile != "" {
		visited[ix.aliasFile] = true
	}

	for _, pattern := range includes {
		expanded := expandPath(pattern)
		var matches []string

		info, err := os.Stat(expanded)
		if err == nil && info.IsDir() {
			// If pattern is a directory, expand to all files inside it
			matches, _ = filepath.Glob(filepath.Join(expanded, "*"))
		} else {
			matches, err = filepath.Glob(expanded)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("glob %q: %w", pattern, err))
				continue
			}
		}

		for _, match := range matches {
			if visited[match] {
				continue
			}

			info, err := os.Stat(match)
			if err != nil || info.IsDir() {
				continue
			}

			if isExcluded(match, excludes) {
				continue
			}

			visited[match] = true
			if err := ix.indexFile(match, false, &result); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("index %s: %w", match, err))
			}
		}
	}

	result.Duration = time.Since(start)
	return result
}

func isExcluded(path string, excludes []string) bool {
	for _, ex := range excludes {
		expandedEx := expandPath(ex)
		if strings.Contains(expandedEx, "/") || strings.Contains(expandedEx, string(filepath.Separator)) {
			if matched, _ := filepath.Match(expandedEx, path); matched || path == expandedEx {
				return true
			}
		} else {
			if matched, _ := filepath.Match(expandedEx, filepath.Base(path)); matched || filepath.Base(path) == expandedEx {
				return true
			}
		}
	}
	return false
}

// NeedsRefresh reports whether the SQLite index needs to be updated.
// It returns true if:
// 1. cacheInterval <= 0 or index is empty
// 2. time.Since(lastIndexed) > cacheInterval
// 3. Primary alias file or any index_folders file is missing from DB or has an updated mtime.
func (ix *Indexer) NeedsRefresh() bool {
	if ix.cacheInterval <= 0 {
		return true
	}

	var lastIndexed string
	err := ix.store.db.QueryRow(
		"SELECT MAX(indexed_at) FROM source_files",
	).Scan(&lastIndexed)
	if err != nil || lastIndexed == "" {
		return true // never indexed
	}

	t, err := time.Parse(time.RFC3339, lastIndexed)
	if err == nil && time.Since(t) > time.Duration(ix.cacheInterval)*time.Second {
		return true // cache interval expired
	}

	// Check if primary alias file needs indexing
	if ix.aliasFile != "" && ix.fileNeedsIndex(ix.aliasFile) {
		return true
	}

	// Check if any matching file in index_folders needs indexing
	var includes []string
	var excludes []string
	for _, raw := range ix.indexFolders {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "!") {
			excludes = append(excludes, strings.TrimPrefix(trimmed, "!"))
		} else {
			includes = append(includes, trimmed)
		}
	}

	visited := make(map[string]bool)
	if ix.aliasFile != "" {
		visited[ix.aliasFile] = true
	}

	for _, pattern := range includes {
		expanded := expandPath(pattern)
		var matches []string

		info, err := os.Stat(expanded)
		if err == nil && info.IsDir() {
			matches, _ = filepath.Glob(filepath.Join(expanded, "*"))
		} else {
			matches, _ = filepath.Glob(expanded)
		}

		for _, match := range matches {
			if visited[match] {
				continue
			}
			info, err := os.Stat(match)
			if err != nil || info.IsDir() {
				continue
			}
			if isExcluded(match, excludes) {
				continue
			}
			visited[match] = true
			if ix.fileNeedsIndex(match) {
				return true
			}
		}
	}

	return false
}

func (ix *Indexer) fileNeedsIndex(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	fileMtime := info.ModTime().Unix()
	storedMtime, err := ix.store.SourceMtime(path)
	if err != nil || storedMtime == 0 {
		return true // unindexed or query error
	}
	return storedMtime != fileMtime
}

// indexFile reads a single shell file and upserts its aliases into the store.
// It skips re-indexing if the file mtime has not changed.
func (ix *Indexer) indexFile(path string, global bool, result *IndexResult) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	fileMtime := info.ModTime().Unix()

	// Check if mtime is unchanged — skip if so.
	storedMtime, err := ix.store.SourceMtime(path)
	if err != nil {
		return err
	}
	if storedMtime == fileMtime {
		result.FilesSkipped++
		return nil
	}

	// Use a temporary alias.Manager to parse the shell file format.
	mgr := alias.NewManager(path, ix.shell, "", "", nil)
	aliases, err := mgr.Load()
	if err != nil {
		return err
	}

	// Set source on each alias.
	for i := range aliases {
		aliases[i].SourceFile = path
	}

	if err := ix.store.BulkUpsert(aliases, path, fileMtime, global); err != nil {
		return err
	}

	result.FilesScanned++
	result.AliasesStored += len(aliases)
	return nil
}

func expandPath(value string) string {
	expanded := os.ExpandEnv(value)
	if expanded == "" {
		return expanded
	}
	if expanded == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
		return expanded
	}
	if strings.HasPrefix(expanded, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(expanded, "~/"))
		}
	}
	return expanded
}
