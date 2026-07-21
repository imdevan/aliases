package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aliases/internal/domain"
)

func tempDB(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test-index.db")
	s, err := NewStoreAt(dbPath)
	if err != nil {
		t.Fatalf("NewStoreAt: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestStore_UpsertAndSearch(t *testing.T) {
	s := tempDB(t)

	a := domain.Alias{
		Name:        "gs",
		Value:       "git status",
		Description: "show git status",
		SourceFile:  "/home/user/.aliases/aliases.zsh",
	}
	if err := s.UpsertAlias(a, 1000, true); err != nil {
		t.Fatalf("UpsertAlias: %v", err)
	}

	// Search by name.
	results, err := s.Search("gs")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "gs" || results[0].Value != "git status" {
		t.Errorf("unexpected result: %+v", results[0])
	}

	// Search by value.
	results, err = s.Search("git")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	// Search by description.
	results, err = s.Search("show")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	// No results.
	results, err = s.Search("nonexistent")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestStore_BulkUpsert(t *testing.T) {
	s := tempDB(t)

	source := "/home/user/.aliases/aliases.zsh"
	aliases := []domain.Alias{
		{Name: "gs", Value: "git status", Description: "git status"},
		{Name: "ga", Value: "git add", Description: "git add"},
		{Name: "gc", Value: "git commit", Description: "git commit"},
	}

	if err := s.BulkUpsert(aliases, source, 1000, true); err != nil {
		t.Fatalf("BulkUpsert: %v", err)
	}

	all, err := s.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 aliases, got %d", len(all))
	}

	// Re-upsert with one removed — stale entry should be deleted.
	aliases2 := []domain.Alias{
		{Name: "gs", Value: "git status -s", Description: "short status"},
		{Name: "ga", Value: "git add", Description: "git add"},
	}

	if err := s.BulkUpsert(aliases2, source, 2000, true); err != nil {
		t.Fatalf("BulkUpsert: %v", err)
	}

	all, err = s.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 aliases after re-upsert, got %d", len(all))
	}

	// Check the updated value.
	results, err := s.Search("gs")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 || results[0].Value != "git status -s" {
		t.Errorf("expected updated value, got %+v", results)
	}
}

func TestStore_BySource(t *testing.T) {
	s := tempDB(t)

	src1 := "/home/user/.aliases/aliases.zsh"
	src2 := "/home/user/.dotfiles/git.zsh"

	s.BulkUpsert([]domain.Alias{
		{Name: "gs", Value: "git status"},
	}, src1, 1000, true)

	s.BulkUpsert([]domain.Alias{
		{Name: "ll", Value: "ls -la"},
	}, src2, 1000, false)

	r1, _ := s.BySource(src1)
	if len(r1) != 1 || r1[0].Name != "gs" {
		t.Errorf("BySource src1: %+v", r1)
	}

	r2, _ := s.BySource(src2)
	if len(r2) != 1 || r2[0].Name != "ll" {
		t.Errorf("BySource src2: %+v", r2)
	}
}

func TestStore_DeleteBySource(t *testing.T) {
	s := tempDB(t)

	source := "/home/user/.aliases/aliases.zsh"
	s.BulkUpsert([]domain.Alias{
		{Name: "gs", Value: "git status"},
		{Name: "ga", Value: "git add"},
	}, source, 1000, true)

	if err := s.DeleteBySource(source); err != nil {
		t.Fatalf("DeleteBySource: %v", err)
	}

	all, _ := s.All()
	if len(all) != 0 {
		t.Fatalf("expected 0 aliases after delete, got %d", len(all))
	}

	// Source mtime should be gone too.
	mtime, _ := s.SourceMtime(source)
	if mtime != 0 {
		t.Errorf("expected mtime=0 after delete, got %d", mtime)
	}
}

func TestStore_SourceMtime(t *testing.T) {
	s := tempDB(t)

	source := "/home/user/.aliases/aliases.zsh"

	// Not indexed yet.
	mtime, err := s.SourceMtime(source)
	if err != nil {
		t.Fatalf("SourceMtime: %v", err)
	}
	if mtime != 0 {
		t.Errorf("expected 0, got %d", mtime)
	}

	// After indexing.
	s.BulkUpsert([]domain.Alias{
		{Name: "gs", Value: "git status"},
	}, source, 12345, true)

	mtime, err = s.SourceMtime(source)
	if err != nil {
		t.Fatalf("SourceMtime: %v", err)
	}
	if mtime != 12345 {
		t.Errorf("expected 12345, got %d", mtime)
	}
}

func TestStore_Stats(t *testing.T) {
	s := tempDB(t)

	s.BulkUpsert([]domain.Alias{
		{Name: "gs", Value: "git status"},
		{Name: "ga", Value: "git add"},
	}, "/src1", 1000, true)

	s.BulkUpsert([]domain.Alias{
		{Name: "ll", Value: "ls -la"},
	}, "/src2", 1000, false)

	totalAliases, totalSources, err := s.Stats()
	if err != nil {
		t.Fatalf("Stats: %v", err)
	}
	if totalAliases != 3 {
		t.Errorf("expected 3 aliases, got %d", totalAliases)
	}
	if totalSources != 2 {
		t.Errorf("expected 2 sources, got %d", totalSources)
	}
}

func TestIndexer_Refresh(t *testing.T) {
	s := tempDB(t)

	// Create a temp alias file.
	dir := t.TempDir()
	aliasFile := filepath.Join(dir, "aliases.zsh")
	content := `alias gs="git status" # show git status
alias ga="git add"
alias gc="git commit -m"
`
	if err := os.WriteFile(aliasFile, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg := domain.Config{
		AliasFile:     aliasFile,
		Shell:         "zsh",
		IndexFolders:  []string{},
		CacheInterval: 300,
	}

	ix := NewIndexer(s, cfg)
	result := ix.Refresh()

	if len(result.Errors) > 0 {
		t.Fatalf("Refresh errors: %v", result.Errors)
	}
	if result.FilesScanned != 1 {
		t.Errorf("expected 1 file scanned, got %d", result.FilesScanned)
	}
	if result.AliasesStored != 3 {
		t.Errorf("expected 3 aliases stored, got %d", result.AliasesStored)
	}

	// Verify in DB.
	all, _ := s.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 aliases in DB, got %d", len(all))
	}

	// Second refresh should skip (mtime unchanged).
	result2 := ix.Refresh()
	if result2.FilesSkipped != 1 {
		t.Errorf("expected 1 file skipped, got %d", result2.FilesSkipped)
	}
	if result2.FilesScanned != 0 {
		t.Errorf("expected 0 files scanned on second run, got %d", result2.FilesScanned)
	}
}

func TestIndexer_RefreshWithIndexFolders(t *testing.T) {
	s := tempDB(t)

	dir := t.TempDir()

	// Primary alias file.
	aliasFile := filepath.Join(dir, "aliases.zsh")
	os.WriteFile(aliasFile, []byte(`alias gs="git status"`), 0o644)

	// Extra indexed file.
	extraDir := filepath.Join(dir, "extra")
	os.MkdirAll(extraDir, 0o755)
	extraFile := filepath.Join(extraDir, "docker.zsh")
	os.WriteFile(extraFile, []byte(`alias dps="docker ps"
alias dimg="docker images"`), 0o644)

	cfg := domain.Config{
		AliasFile:     aliasFile,
		Shell:         "zsh",
		IndexFolders:  []string{filepath.Join(extraDir, "*.zsh")},
		CacheInterval: 300,
	}

	ix := NewIndexer(s, cfg)
	result := ix.Refresh()

	if len(result.Errors) > 0 {
		t.Fatalf("Refresh errors: %v", result.Errors)
	}
	if result.FilesScanned != 2 {
		t.Errorf("expected 2 files scanned, got %d", result.FilesScanned)
	}
	if result.AliasesStored != 3 {
		t.Errorf("expected 3 aliases stored, got %d", result.AliasesStored)
	}

	all, _ := s.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 aliases in DB, got %d", len(all))
	}
}

func TestIndexer_NeedsRefresh(t *testing.T) {
	s := tempDB(t)

	cfg := domain.Config{
		AliasFile:     "/nonexistent",
		Shell:         "zsh",
		CacheInterval: 300,
	}

	ix := NewIndexer(s, cfg)

	// Before any indexing — should need refresh.
	if !ix.NeedsRefresh() {
		t.Error("expected NeedsRefresh=true before first index")
	}

	// Create a real file and index it.
	dir := t.TempDir()
	aliasFile := filepath.Join(dir, "aliases.zsh")
	os.WriteFile(aliasFile, []byte(`alias gs="git status"`), 0o644)

	cfg.AliasFile = aliasFile
	ix = NewIndexer(s, cfg)
	ix.Refresh()

	// Right after indexing — should NOT need refresh.
	if ix.NeedsRefresh() {
		t.Error("expected NeedsRefresh=false right after indexing")
	}

	// Adding a new folder/file to IndexFolders should trigger NeedsRefresh=true immediately
	extraFile := filepath.Join(dir, "extra.zsh")
	os.WriteFile(extraFile, []byte(`alias ex="echo extra"`), 0o644)
	cfg.IndexFolders = []string{extraFile}
	ix = NewIndexer(s, cfg)

	if !ix.NeedsRefresh() {
		t.Error("expected NeedsRefresh=true after adding unindexed file to IndexFolders")
	}

	// With cache_interval=0 — should always need refresh.
	cfg.CacheInterval = 0
	ix = NewIndexer(s, cfg)
	if !ix.NeedsRefresh() {
		t.Error("expected NeedsRefresh=true with cache_interval=0")
	}
}

func TestIndexer_ExclusionPatternsAndDirectoryScanning(t *testing.T) {
	s := tempDB(t)

	dir := t.TempDir()

	aliasFile := filepath.Join(dir, "primary.zsh")
	os.WriteFile(aliasFile, []byte(`alias p1="echo primary"`), 0o644)

	extraDir := filepath.Join(dir, "extra")
	os.MkdirAll(extraDir, 0o755)

	f1 := filepath.Join(extraDir, "git.zsh")
	os.WriteFile(f1, []byte(`alias g1="git status"`), 0o644)

	f2 := filepath.Join(extraDir, "secret.zsh")
	os.WriteFile(f2, []byte(`alias s1="echo secret"`), 0o644)

	f3 := filepath.Join(extraDir, "backup.tmp")
	os.WriteFile(f3, []byte(`alias b1="echo temp"`), 0o644)

	cfg := domain.Config{
		AliasFile:    aliasFile,
		Shell:        "zsh",
		IndexFolders: []string{extraDir, "!" + f2, "!*.tmp"},
	}

	ix := NewIndexer(s, cfg)
	res := ix.Refresh()

	if res.FilesScanned != 2 { // primary + extra/git.zsh (secret.zsh and backup.tmp excluded)
		t.Errorf("expected 2 files scanned, got %d", res.FilesScanned)
	}

	all, err := s.All()
	if err != nil {
		t.Fatalf("s.All(): %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(all))
	}
	names := map[string]bool{}
	for _, a := range all {
		names[a.Name] = true
	}
	if !names["p1"] || !names["g1"] {
		t.Errorf("unexpected aliases: %+v", all)
	}
}

func TestStore_GlobalAndSuffixSchema(t *testing.T) {
	s := tempDB(t)

	source := "/path/to/file.zsh"
	a1 := domain.Alias{Name: "gs", Value: "git status"}
	a2 := domain.Alias{Name: "gc-", Value: "git commit"}

	if err := s.BulkUpsert([]domain.Alias{a1, a2}, source, 100, true); err != nil {
		t.Fatalf("BulkUpsert: %v", err)
	}

	all, err := s.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(all))
	}

	for _, a := range all {
		if !a.Global {
			t.Errorf("expected Global=true for %s", a.Name)
		}
		if a.Name == "gc-" && !a.Suffix {
			t.Errorf("expected Suffix=true for %s", a.Name)
		}
		if a.Name == "gs" && a.Suffix {
			t.Errorf("expected Suffix=false for %s", a.Name)
		}
	}
}

func TestIndexer_RefreshAsyncInProcess(t *testing.T) {
	s := tempDB(t)

	dir := t.TempDir()
	aliasFile := filepath.Join(dir, "aliases.zsh")
	os.WriteFile(aliasFile, []byte(`alias gs="git status"`), 0o644)

	cfg := domain.Config{
		AliasFile:     aliasFile,
		Shell:         "zsh",
		CacheInterval: 300,
	}

	ix := NewIndexer(s, cfg)
	if err := ix.RefreshAsync(false); err != nil {
		t.Fatalf("RefreshAsync error: %v", err)
	}

	// Give goroutine a moment to complete
	for i := 0; i < 50; i++ {
		if !ix.NeedsRefresh() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if ix.NeedsRefresh() {
		t.Error("expected NeedsRefresh=false after RefreshAsync")
	}
}
