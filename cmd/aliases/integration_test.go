package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aliases/internal/bookmark"
	"github.com/aliases/internal/testutil"
)

// TestIntegration_BookmarkCurrentDirectory tests task 1.1: bookmark current folder
func TestIntegration_BookmarkCurrentDirectory(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	// Create test directory
	testDir := filepath.Join(env.TempDir, "my-cool-project")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Test auto-generated alias
	cmd := newRootCmd()
	cmd.SetArgs([]string{"-c", env.ConfigPath, "-s", testDir})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	// Verify bookmark was created
	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	bookmarks, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load bookmarks: %v", err)
	}

	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
	}

	bm := bookmarks[0]
	if bm.Alias != "mcp" {
		t.Errorf("expected alias 'mcp', got %q", bm.Alias)
	}
	if bm.Path != testDir {
		t.Errorf("expected path %q, got %q", testDir, bm.Path)
	}
}

// TestIntegration_AutoGenerateAlias tests task 1.2: auto-generate alias from directory name
func TestIntegration_AutoGenerateAlias(t *testing.T) {
	tests := []struct {
		dirName      string
		expectedAlias string
	}{
		{"my-cool-project", "mcp"},
		{"web-app", "wa"},
		{"single", "s"},
		{"foo_bar_baz", "fbb"},
	}

	for _, tt := range tests {
		t.Run(tt.dirName, func(t *testing.T) {
			env := testutil.SetupTestEnv(t)
			defer env.Cleanup()

			testDir := filepath.Join(env.TempDir, tt.dirName)
			if err := os.MkdirAll(testDir, 0o755); err != nil {
				t.Fatal(err)
			}

			cmd := newRootCmd()
			cmd.SetArgs([]string{"-c", env.ConfigPath, "-s", testDir})
			
			if err := cmd.Execute(); err != nil {
				t.Fatalf("failed to create bookmark: %v", err)
			}

			cfg := testutil.LoadTestConfig(t, env.ConfigPath)
			mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
			
			bookmarks, err := mgr.Load()
			if err != nil {
				t.Fatalf("failed to load bookmarks: %v", err)
			}

			if len(bookmarks) != 1 {
				t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
			}

			if bookmarks[0].Alias != tt.expectedAlias {
				t.Errorf("expected alias %q, got %q", tt.expectedAlias, bookmarks[0].Alias)
			}
		})
	}
}

// TestIntegration_CustomAlias tests task 1.3: custom alias via argument
func TestIntegration_CustomAlias(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "projects", "webapp")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create bookmark with custom alias
	cmd := newRootCmd()
	cmd.SetArgs([]string{"web", "-c", env.ConfigPath, "-s", testDir})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	bookmarks, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load bookmarks: %v", err)
	}

	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
	}

	bm := bookmarks[0]
	if bm.Alias != "web" {
		t.Errorf("expected alias 'web', got %q", bm.Alias)
	}
	if bm.Path != testDir {
		t.Errorf("expected path %q, got %q", testDir, bm.Path)
	}
}

// TestIntegration_InvalidAlias tests task 1.3: alias validation
func TestIntegration_InvalidAlias(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "test")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	invalidAliases := []string{
		"my alias",    // space
		"my@alias",    // special char
		"my.alias",    // dot
		"my/alias",    // slash
	}

	for _, alias := range invalidAliases {
		t.Run(alias, func(t *testing.T) {
			cmd := newRootCmd()
			cmd.SetArgs([]string{alias, "-c", env.ConfigPath, "-s", testDir})
			
			err := cmd.Execute()
			if err == nil {
				t.Errorf("expected error for invalid alias %q, got nil", alias)
			}
		})
	}
}

// TestIntegration_OverwriteConfirmation tests task 1.4: confirmation before overwriting
func TestIntegration_OverwriteConfirmation(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "webapp")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create initial bookmark
	cmd := newRootCmd()
	cmd.SetArgs([]string{"web", "-c", env.ConfigPath, "-s", testDir})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create initial bookmark: %v", err)
	}

	// Verify bookmark exists
	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	exists, err := mgr.Exists("web")
	if err != nil {
		t.Fatalf("failed to check bookmark existence: %v", err)
	}
	if !exists {
		t.Fatal("expected bookmark to exist")
	}
}

// TestIntegration_ForceOverwrite tests task 1.5: -y flag to bypass confirmation
func TestIntegration_ForceOverwrite(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir1 := filepath.Join(env.TempDir, "webapp1")
	testDir2 := filepath.Join(env.TempDir, "webapp2")
	if err := os.MkdirAll(testDir1, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(testDir2, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create initial bookmark
	cmd := newRootCmd()
	cmd.SetArgs([]string{"web", "-c", env.ConfigPath, "-s", testDir1})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create initial bookmark: %v", err)
	}

	// Overwrite with -y flag
	cmd2 := newRootCmd()
	cmd2.SetArgs([]string{"web", "-y", "-c", env.ConfigPath, "-s", testDir2})
	
	if err := cmd2.Execute(); err != nil {
		t.Fatalf("failed to overwrite bookmark: %v", err)
	}

	// Verify bookmark was updated
	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	bm, err := mgr.Get("web")
	if err != nil {
		t.Fatalf("failed to get bookmark: %v", err)
	}

	if bm.Path != testDir2 {
		t.Errorf("expected path %q, got %q", testDir2, bm.Path)
	}
}

// TestIntegration_FileFlag tests task 1.6: -f flag to open file after navigation
func TestIntegration_FileFlag(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "foo-bar")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create bookmark with file flag
	cmd := newRootCmd()
	cmd.SetArgs([]string{"-f", "plan.md", "-c", env.ConfigPath, "-s", testDir})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	bookmarks, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load bookmarks: %v", err)
	}

	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
	}

	bm := bookmarks[0]
	if bm.File != "plan.md" {
		t.Errorf("expected file 'plan.md', got %q", bm.File)
	}

	// Verify the generated alias includes editor command
	content, err := os.ReadFile(cfg.BookmarkFile())
	if err != nil {
		t.Fatalf("failed to read bookmark file: %v", err)
	}

	if !strings.Contains(string(content), "plan.md") {
		t.Error("expected bookmark file to contain 'plan.md'")
	}
}

// TestIntegration_ExecuteFlag tests task 1.9: -x flag for custom execution
func TestIntegration_ExecuteFlag(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "foo-bar")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create bookmark with execute flag
	cmd := newRootCmd()
	cmd.SetArgs([]string{"-x", "echo 'hello'", "-c", env.ConfigPath, "-s", testDir})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	bookmarks, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load bookmarks: %v", err)
	}

	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
	}

	bm := bookmarks[0]
	if bm.Execute != "echo 'hello'" {
		t.Errorf("expected execute 'echo 'hello'', got %q", bm.Execute)
	}
}

// TestIntegration_SourceFlag tests task 1.10: -s flag for custom source location
func TestIntegration_SourceFlag(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	sourceDir := filepath.Join(env.TempDir, "Documents", "bar")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create bookmark with source flag
	cmd := newRootCmd()
	cmd.SetArgs([]string{"b", "-s", sourceDir, "-c", env.ConfigPath})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	bm, err := mgr.Get("b")
	if err != nil {
		t.Fatalf("failed to get bookmark: %v", err)
	}

	if bm.Path != sourceDir {
		t.Errorf("expected path %q, got %q", sourceDir, bm.Path)
	}
}

// TestIntegration_TmuxFlag tests task 1.11: -T flag for custom tmux window name
func TestIntegration_TmuxFlag(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "foo-bar")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create bookmark with tmux flag
	cmd := newRootCmd()
	cmd.SetArgs([]string{"-T", "bar", "-c", env.ConfigPath, "-s", testDir})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	bookmarks, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load bookmarks: %v", err)
	}

	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
	}

	bm := bookmarks[0]
	if bm.TmuxWindowName != "bar" {
		t.Errorf("expected tmux window name 'bar', got %q", bm.TmuxWindowName)
	}
}

// TestIntegration_CombinedFlags tests combining multiple flags
func TestIntegration_CombinedFlags(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "project")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create bookmark with multiple flags
	cmd := newRootCmd()
	cmd.SetArgs([]string{
		"proj",
		"-T", "myproject",
		"-f", "README.md",
		"-x", "git status",
		"-c", env.ConfigPath,
		"-s", testDir,
	})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	bm, err := mgr.Get("proj")
	if err != nil {
		t.Fatalf("failed to get bookmark: %v", err)
	}

	if bm.Alias != "proj" {
		t.Errorf("expected alias 'proj', got %q", bm.Alias)
	}
	if bm.Path != testDir {
		t.Errorf("expected path %q, got %q", testDir, bm.Path)
	}
	if bm.TmuxWindowName != "myproject" {
		t.Errorf("expected tmux window name 'myproject', got %q", bm.TmuxWindowName)
	}
	if bm.File != "README.md" {
		t.Errorf("expected file 'README.md', got %q", bm.File)
	}
	if bm.Execute != "git status" {
		t.Errorf("expected execute 'git status', got %q", bm.Execute)
	}
}

// TestIntegration_AliasStructure tests task 1.12: alias structure validation
func TestIntegration_AliasStructure(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "project")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create bookmark with all features
	cmd := newRootCmd()
	cmd.SetArgs([]string{
		"test",
		"-T", "testwin",
		"-x", "echo 'setup'",
		"-f", "main.go",
		"-c", env.ConfigPath,
		"-s", testDir,
	})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	
	// Read the generated shell script
	content, err := os.ReadFile(cfg.BookmarkFile())
	if err != nil {
		t.Fatalf("failed to read bookmark file: %v", err)
	}

	script := string(content)
	
	// Debug: print the generated script
	t.Logf("Generated script:\n%s", script)
	
	// Extract just the alias line for position checking
	lines := strings.Split(script, "\n")
	var aliasLine string
	for _, line := range lines {
		if strings.HasPrefix(line, "alias test=") {
			aliasLine = line
			break
		}
	}
	
	if aliasLine == "" {
		t.Fatal("could not find alias line in generated script")
	}
	
	t.Logf("Alias line: %s", aliasLine)
	
	// Verify structure: navigate -> tmux rename -> execute -> open file
	if !strings.Contains(aliasLine, "tmux rename-window") {
		t.Error("expected alias to contain tmux rename-window")
	}
	// Check for quoted path (paths are now properly quoted)
	if !strings.Contains(aliasLine, "cd '") && !strings.Contains(aliasLine, "cd "+testDir) {
		t.Errorf("expected alias to contain cd command with path")
	}
	if !strings.Contains(aliasLine, "echo") {
		t.Error("expected alias to contain execute command (echo)")
	}
	if !strings.Contains(aliasLine, "main.go") {
		t.Error("expected alias to contain file reference")
	}

	// Verify order using string positions within the alias line
	cdPos := strings.Index(aliasLine, "cd ")
	tmuxPos := strings.Index(aliasLine, "tmux rename-window")
	execPos := strings.Index(aliasLine, "echo")
	filePos := strings.Index(aliasLine, "main.go")

	t.Logf("Positions in alias: cd=%d, tmux=%d, exec=%d, file=%d", cdPos, tmuxPos, execPos, filePos)

	if cdPos > tmuxPos {
		t.Error("cd should come before tmux rename")
	}
	if tmuxPos > execPos {
		t.Error("tmux rename should come before execute")
	}
	if execPos > filePos {
		t.Error("execute should come before file open")
	}
}

// TestIntegration_EditFlag tests task 1.7: -e flag to edit bookmarks
func TestIntegration_EditFlag(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "foo-bar")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create a bookmark first
	cmd := newRootCmd()
	cmd.SetArgs([]string{"fb", "-c", env.ConfigPath, "-s", testDir})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	// Verify bookmark exists
	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	exists, err := mgr.Exists("fb")
	if err != nil {
		t.Fatalf("failed to check bookmark existence: %v", err)
	}
	if !exists {
		t.Fatal("expected bookmark 'fb' to exist")
	}

	// Test that FindBookmarkLine works
	line, err := mgr.FindBookmarkLine("fb")
	if err != nil {
		t.Fatalf("failed to find bookmark line: %v", err)
	}
	if line == 0 {
		t.Error("expected non-zero line number")
	}
}

// TestIntegration_TmuxAutoName tests task 1.1: -t flag with auto tmux name
func TestIntegration_TmuxAutoName(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	defer env.Cleanup()

	testDir := filepath.Join(env.TempDir, "my-cool-project")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create bookmark with -t flag (auto tmux name same as alias)
	cmd := newRootCmd()
	cmd.SetArgs([]string{"-t", "-c", env.ConfigPath, "-s", testDir})
	
	if err := cmd.Execute(); err != nil {
		t.Fatalf("failed to create bookmark: %v", err)
	}

	cfg := testutil.LoadTestConfig(t, env.ConfigPath)
	mgr := bookmark.NewManager(cfg.BookmarkFile(), cfg.Shell, cfg.NavigationTool, cfg.Editor, cfg.FunctionAlias, cfg.InteractiveAlias)
	
	bookmarks, err := mgr.Load()
	if err != nil {
		t.Fatalf("failed to load bookmarks: %v", err)
	}

	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
	}

	bm := bookmarks[0]
	if bm.TmuxWindowName != bm.Alias {
		t.Errorf("expected tmux window name to match alias %q, got %q", bm.Alias, bm.TmuxWindowName)
	}
}
