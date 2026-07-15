package bookmark

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bookmark/internal/domain"
)

func TestManager_BuildNavigationCommand_WithQuotes(t *testing.T) {
	tests := []struct {
		name     string
		shell    string
		bookmark domain.Bookmark
		want     string
	}{
		{
			name:  "bash single quote in path",
			shell: "bash",
			bookmark: domain.Bookmark{
				Alias: "test",
				Path:  "/home/user/let's-go/test",
			},
			want: "cd '/home/user/let'\\''s-go/test'",
		},
		{
			name:  "fish single quote in path",
			shell: "fish",
			bookmark: domain.Bookmark{
				Alias: "test",
				Path:  "/home/user/let's-go/test",
			},
			want: "cd '/home/user/let\\'s-go/test'",
		},
		{
			name:  "nushell single quote in path",
			shell: "nu",
			bookmark: domain.Bookmark{
				Alias: "test",
				Path:  "/home/user/let's-go/test",
			},
			want: "cd '/home/user/let''s-go/test'",
		},
		{
			name:  "bash single quote in file and tmux",
			shell: "bash",
			bookmark: domain.Bookmark{
				Alias:          "test",
				Path:           "/home/user/project",
				File:           "some'file.txt",
				TmuxWindowName: "tmux'win",
			},
			want: "cd '/home/user/project' && tmux rename-window 'tmux'\\''win' && nvim 'some'\\''file.txt'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager("/tmp/bookmarks.sh", tt.shell, "cd", "nvim", "false", "false")
			got := m.BuildNavigationCommand(tt.bookmark)
			if got != tt.want {
				t.Errorf("BuildNavigationCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestManager_SaveAndLoad_WithPipesAndNewlines(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bookmarks.sh")

	m := NewManager(filePath, "bash", "cd", "nvim", "false", "false")

	bookmarks := []domain.Bookmark{
		{
			Alias:          "test1",
			Path:           "/home/user/path|with|pipes",
			Description:    "Description with | pipe and \n newline",
			TmuxWindowName: "tmux|pipe",
			Execute:        "echo 'hello' | grep 'h'",
			PostJumpScript: "echo 'post' | wc -l",
			File:           "file|pipe.txt",
			CreatedAt:      time.Now().Truncate(time.Second),
			UpdatedAt:      time.Now().Truncate(time.Second),
		},
		{
			Alias:          "test2",
			Path:           "/home/user/path\\with\\backslashes",
			Description:    "Description with \\ backslash",
			TmuxWindowName: "tmux\\bs",
			Execute:        "echo \\\"hello\\\"",
			PostJumpScript: "echo \\\"post\\\"",
			File:           "file\\bs.txt",
			CreatedAt:      time.Now().Truncate(time.Second),
			UpdatedAt:      time.Now().Truncate(time.Second),
		},
	}

	for _, bm := range bookmarks {
		if err := m.Add(bm); err != nil {
			t.Fatalf("Add() error = %v", err)
		}
	}

	// Verify the file was written
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}
	contentStr := string(content)

	// Verify no unescaped pipe characters in metadata lines
	lines := strings.Split(contentStr, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# BM:") {
			// Find how many unescaped pipes there are.
			// Since there should be exactly 8 unescaped pipes to separate 9 fields,
			// let's check our splitMetadata function.
			parts := splitMetadata(strings.TrimPrefix(line, "# BM: "))
			if len(parts) != 9 {
				t.Errorf("expected 9 fields in metadata, got %d from line:\n%s", len(parts), line)
			}
		}
	}

	// Load bookmarks back
	loaded, err := m.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != len(bookmarks) {
		t.Fatalf("expected %d bookmarks, got %d", len(bookmarks), len(loaded))
	}

	for i, want := range bookmarks {
		got := loaded[i]
		if got.Alias != want.Alias {
			t.Errorf("Bookmark[%d].Alias = %q, want %q", i, got.Alias, want.Alias)
		}
		if got.Path != want.Path {
			t.Errorf("Bookmark[%d].Path = %q, want %q", i, got.Path, want.Path)
		}
		if got.Description != want.Description {
			t.Errorf("Bookmark[%d].Description = %q, want %q", i, got.Description, want.Description)
		}
		if got.TmuxWindowName != want.TmuxWindowName {
			t.Errorf("Bookmark[%d].TmuxWindowName = %q, want %q", i, got.TmuxWindowName, want.TmuxWindowName)
		}
		if got.Execute != want.Execute {
			t.Errorf("Bookmark[%d].Execute = %q, want %q", i, got.Execute, want.Execute)
		}
		if got.PostJumpScript != want.PostJumpScript {
			t.Errorf("Bookmark[%d].PostJumpScript = %q, want %q", i, got.PostJumpScript, want.PostJumpScript)
		}
		if got.File != want.File {
			t.Errorf("Bookmark[%d].File = %q, want %q", i, got.File, want.File)
		}
	}
}
