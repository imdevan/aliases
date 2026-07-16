package bookmark

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aliases/internal/domain"
)

func TestManager_InteractiveAliasGeneration(t *testing.T) {
	tests := []struct {
		name              string
		shell             string
		interactiveAlias  string
		wantWrapper       bool
		wantFunction      string
	}{
		{
			name:             "bash with default interactive alias",
			shell:            "bash",
			interactiveAlias: "bm",
			wantWrapper:      true,
			wantFunction:     "bm() {",
		},
		{
			name:             "bash with custom interactive alias",
			shell:            "bash",
			interactiveAlias: "goto",
			wantWrapper:      true,
			wantFunction:     "goto() {",
		},
		{
			name:             "bash with disabled interactive alias",
			shell:            "bash",
			interactiveAlias: "false",
			wantWrapper:      false,
			wantFunction:     "",
		},
		{
			name:             "fish with default interactive alias",
			shell:            "fish",
			interactiveAlias: "bm",
			wantWrapper:      true,
			wantFunction:     "function bm",
		},
		{
			name:             "zsh with default interactive alias",
			shell:            "zsh",
			interactiveAlias: "bm",
			wantWrapper:      true,
			wantFunction:     "bm() {",
		},
		{
			name:             "nushell with interactive alias",
			shell:            "nu",
			interactiveAlias: "bm",
			wantWrapper:      true,
			wantFunction:     "def bm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "bookmarks.sh")

			m := NewManager(filePath, tt.shell, "cd", "nvim", "false", tt.interactiveAlias)

			// Add a test bookmark
			bookmark := domain.Bookmark{
				Alias:       "test",
				Path:        "/tmp/test",
				Description: "Test bookmark",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			if err := m.Add(bookmark); err != nil {
				t.Fatalf("Add() error = %v", err)
			}

			// Read generated file
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read generated file: %v", err)
			}

			contentStr := string(content)

			// Check for interactive wrapper presence
			hasWrapper := strings.Contains(contentStr, "Interactive bookmark navigation function")
			if hasWrapper != tt.wantWrapper {
				t.Errorf("interactive wrapper presence = %v, want %v", hasWrapper, tt.wantWrapper)
			}

			// Check for specific function declaration
			if tt.wantFunction != "" {
				if !strings.Contains(contentStr, tt.wantFunction) {
					t.Errorf("expected function declaration %q not found in:\n%s", tt.wantFunction, contentStr)
				}
			}

			// Verify bookmark alias is still present
			if !strings.Contains(contentStr, "test") {
				t.Error("bookmark alias not found in generated file")
			}

			// Verify interactive wrapper calls bookmark -i
			if tt.wantWrapper {
				if !strings.Contains(contentStr, "bookmark -i") {
					t.Error("expected 'bookmark -i' in interactive wrapper")
				}
				
				// Verify CLICOLOR_FORCE is set for color support
				if tt.shell == "bash" || tt.shell == "zsh" {
					if !strings.Contains(contentStr, "CLICOLOR_FORCE=1") {
						t.Error("expected 'CLICOLOR_FORCE=1' in bash/zsh interactive wrapper")
					}
				} else if tt.shell == "fish" {
					if !strings.Contains(contentStr, "set -x CLICOLOR_FORCE 1") {
						t.Error("expected 'set -x CLICOLOR_FORCE 1' in fish interactive wrapper")
					}
				} else if tt.shell == "nu" {
					if !strings.Contains(contentStr, "CLICOLOR_FORCE") {
						t.Error("expected 'CLICOLOR_FORCE' in nushell interactive wrapper")
					}
				}
			}
		})
	}
}

func TestManager_InteractiveAliasEmptyString(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bookmarks.sh")

	m := NewManager(filePath, "bash", "cd", "nvim", "false", "")

	bookmark := domain.Bookmark{
		Alias:     "test",
		Path:      "/tmp/test",
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

	// Empty string should disable interactive wrapper
	if strings.Contains(contentStr, "Interactive bookmark navigation") {
		t.Error("interactive wrapper should not be generated when interactiveAlias is empty")
	}
}
