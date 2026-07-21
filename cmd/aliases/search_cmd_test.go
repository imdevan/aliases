package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aliases/internal/domain"
	"github.com/aliases/internal/index"
)

func TestSearchCmd_Execution(t *testing.T) {
	tmpDir := t.TempDir()

	// Set XDG environment variables to isolate index DB
	t.Setenv("XDG_DATA_HOME", filepath.Join(tmpDir, "data"))
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpDir, "config"))

	// Create test database & indexer
	dbPath := filepath.Join(tmpDir, "data", "aliases", "index.db")
	store, err := index.NewStoreAt(dbPath)
	if err != nil {
		t.Fatalf("NewStoreAt: %v", err)
	}
	defer store.Close()

	src1 := filepath.Join(tmpDir, "aliases.zsh")
	src2 := filepath.Join(tmpDir, "docker.zsh")

	aliases1 := []domain.Alias{
		{Name: "gs", Value: "git status", Description: "git status overview"},
		{Name: "ga", Value: "git add", Description: "stage files"},
	}
	aliases2 := []domain.Alias{
		{Name: "dps", Value: "docker ps", Description: "list containers"},
	}

	if err := store.BulkUpsert(aliases1, src1, 1000, true); err != nil {
		t.Fatalf("BulkUpsert 1: %v", err)
	}
	if err := store.BulkUpsert(aliases2, src2, 1000, false); err != nil {
		t.Fatalf("BulkUpsert 2: %v", err)
	}

	t.Run("search query match", func(t *testing.T) {
		cmd := newSearchCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"git"})

		cfg := domain.Config{AliasFile: src1}
		if err := searchIndex(cmd, store, cfg, "git"); err != nil {
			t.Fatalf("searchIndex: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "gs") || !strings.Contains(output, "git status") {
			t.Errorf("expected output to contain 'gs' and 'git status', got:\n%s", output)
		}
		if strings.Contains(output, "dps") {
			t.Errorf("did not expect output to contain 'dps', got:\n%s", output)
		}
	})

	t.Run("search all when query empty", func(t *testing.T) {
		cmd := newSearchCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)

		cfg := domain.Config{AliasFile: src1}
		if err := listAllFromIndex(cmd, store, cfg); err != nil {
			t.Fatalf("listAllFromIndex: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "gs") || !strings.Contains(output, "dps") {
			t.Errorf("expected output to contain both 'gs' and 'dps', got:\n%s", output)
		}
	})

	t.Run("search exact name ending with =", func(t *testing.T) {
		cmd := newSearchCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)

		cfg := domain.Config{AliasFile: src1}
		if err := searchIndex(cmd, store, cfg, "gs="); err != nil {
			t.Fatalf("searchIndex: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "gs") || !strings.Contains(output, "git status") {
			t.Errorf("expected output for gs= to contain 'gs', got:\n%s", output)
		}
		if strings.Contains(output, "ga") {
			t.Errorf("did not expect 'ga' in exact search for 'gs=', got:\n%s", output)
		}
	})
}
