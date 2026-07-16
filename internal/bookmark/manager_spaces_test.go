package bookmark

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aliases/internal/domain"
)

func TestManager_PathsWithSpaces(t *testing.T) {
	tests := []struct {
		name         string
		shell        string
		path         string
		file         string
		wantPath     bool
		wantFile     bool
	}{
		{
			name:     "bash with space in path",
			shell:    "bash",
			path:     "/home/user/my projects/test",
			wantPath: true,
		},
		{
			name:     "zsh with space in path",
			shell:    "zsh",
			path:     "/home/user/my projects/test",
			wantPath: true,
		},
		{
			name:     "fish with space in path",
			shell:    "fish",
			path:     "/home/user/my projects/test",
			wantPath: true,
		},
		{
			name:     "nushell with space in path",
			shell:    "nu",
			path:     "/home/user/my projects/test",
			wantPath: true,
		},
		{
			name:     "bash with space in file path",
			shell:    "bash",
			path:     "/home/user/project",
			file:     "my file.txt",
			wantPath: true,
			wantFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "bookmarks.sh")

			m := NewManager(filePath, tt.shell, "cd", "nvim", "false", "false")

			bookmark := domain.Bookmark{
				Alias:     "test",
				Path:      tt.path,
				File:      tt.file,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := m.Add(bookmark); err != nil {
				t.Fatalf("Add() error = %v", err)
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read generated file: %v", err)
			}

			contentStr := string(content)

			// Verify the path appears in the content (may be escaped)
			if tt.wantPath && !strings.Contains(contentStr, tt.path) {
				t.Errorf("expected path %q in generated content, got:\n%s", tt.path, contentStr)
			}

			// Verify the file appears in the content (may be escaped)
			if tt.wantFile && !strings.Contains(contentStr, tt.file) {
				t.Errorf("expected file %q in generated content, got:\n%s", tt.file, contentStr)
			}

			// Verify the path is not used bare (without quotes) in a cd command
			// This would fail if we have "cd /path with spaces" instead of "cd '/path with spaces'"
			if strings.Contains(tt.path, " ") {
				bareUsage := "cd " + tt.path + " "
				if strings.Contains(contentStr, bareUsage) {
					t.Errorf("path with spaces should be quoted, found bare usage: %s", bareUsage)
				}
			}
		})
	}
}

func TestManager_BuildNavigationCommand_WithSpaces(t *testing.T) {
	m := NewManager("/tmp/bookmarks.sh", "bash", "cd", "nvim", "false", "false")

	tests := []struct {
		name     string
		bookmark domain.Bookmark
		want     string
	}{
		{
			name: "path with spaces",
			bookmark: domain.Bookmark{
				Alias: "test",
				Path:  "/home/user/my projects/test",
			},
			want: "cd '/home/user/my projects/test'",
		},
		{
			name: "path and file with spaces",
			bookmark: domain.Bookmark{
				Alias: "test",
				Path:  "/home/user/my projects/test",
				File:  "my file.txt",
			},
			want: "cd '/home/user/my projects/test' && nvim 'my file.txt'",
		},
		{
			name: "path with spaces and tmux",
			bookmark: domain.Bookmark{
				Alias:          "test",
				Path:           "/home/user/my projects/test",
				TmuxWindowName: "mywindow",
			},
			want: "cd '/home/user/my projects/test' && tmux rename-window 'mywindow'",
		},
		{
			name: "path with spaces and execute",
			bookmark: domain.Bookmark{
				Alias:   "test",
				Path:    "/home/user/my projects/test",
				Execute: "ls -la",
			},
			want: "cd '/home/user/my projects/test' && ls -la",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.BuildNavigationCommand(tt.bookmark)
			if got != tt.want {
				t.Errorf("BuildNavigationCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}
