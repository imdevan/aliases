package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aliases/internal/alias"
	"github.com/aliases/internal/testutil"
)

func TestIntegration_AliasCurrentDirectory(t *testing.T) {
	env := testutil.SetupTestEnv(t)

	// Create test directory
	testDir := filepath.Join(env.TempDir, "my-cool-project")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Change current working directory to testDir so the default alias creation uses it
	origDir, err := os.Getwd()
	if err == nil {
		defer os.Chdir(origDir)
		os.Chdir(testDir)
	}

	cmd := newRootCmd()
	cmd.SetArgs([]string{"-c", env.ConfigPath, "mcp"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create alias: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := alias.NewManager(cfg.ResolvedAliasFile(), cfg.Shell, cfg.FunctionAlias, cfg.InteractiveAlias, cfg.IndexFolders)

	aliases, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load aliases: %v", err)
	}

	if len(aliases) != 1 {
		t.Fatalf("expected 1 alias, got %d", len(aliases))
	}

	al := aliases[0]
	if al.Name != "mcp" {
		t.Errorf("expected name 'mcp', got %q", al.Name)
	}
	if al.Value != "cd "+testDir {
		t.Errorf("expected value 'cd %s', got %q", testDir, al.Value)
	}
}

func TestIntegration_CustomAlias(t *testing.T) {
	env := testutil.SetupTestEnv(t)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"add", "sayhello", "echo hello", "say hello description", "-c", env.ConfigPath, "-y"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to add alias: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := alias.NewManager(cfg.ResolvedAliasFile(), cfg.Shell, cfg.FunctionAlias, cfg.InteractiveAlias, cfg.IndexFolders)

	aliases, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load aliases: %v", err)
	}

	if len(aliases) != 1 {
		t.Fatalf("expected 1 alias, got %d", len(aliases))
	}

	al := aliases[0]
	if al.Name != "sayhello" {
		t.Errorf("expected name 'sayhello', got %q", al.Name)
	}
	if al.Value != "echo hello" {
		t.Errorf("expected value 'echo hello', got %q", al.Value)
	}
	if al.Description != "say hello description" {
		t.Errorf("expected description 'say hello description', got %q", al.Description)
	}
}

func TestIntegration_DeleteAlias(t *testing.T) {
	env := testutil.SetupTestEnv(t)

	// Add alias first
	cmd := newRootCmd()
	cmd.SetArgs([]string{"add", "todelete", "echo delete", "to delete description", "-c", env.ConfigPath, "-y"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to add alias: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := alias.NewManager(cfg.ResolvedAliasFile(), cfg.Shell, cfg.FunctionAlias, cfg.InteractiveAlias, cfg.IndexFolders)

	// Verify it exists
	aliases, _ := mgr.Load()
	if len(aliases) != 1 {
		t.Fatalf("expected 1 alias, got %d", len(aliases))
	}

	// Delete it
	cmd2 := newRootCmd()
	cmd2.SetArgs([]string{"delete", "todelete", "-c", env.ConfigPath, "-f"})
	if err := cmd2.Execute(); err != nil {
		t.Fatalf("failed to delete alias: %v", err)
	}

	// Verify it's gone
	aliases, _ = mgr.Load()
	if len(aliases) != 0 {
		t.Errorf("expected 0 aliases, got %d", len(aliases))
	}
}
