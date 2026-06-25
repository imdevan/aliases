package bookmark

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bookmark/internal/domain"
)

func TestManager_FunctionWrapperGeneration(t *testing.T) {
	tests := []struct {
		name          string
		shell         string
		functionAlias string
		wantWrapper   bool
		wantFunction  string
	}{
		{
			name:          "bash with default function name",
			shell:         "bash",
			functionAlias: "true",
			wantWrapper:   true,
			wantFunction:  "bookmark() {",
		},
		{
			name:          "bash with custom function name",
			shell:         "bash",
			functionAlias: "bm",
			wantWrapper:   true,
			wantFunction:  "bm() {",
		},
		{
			name:          "bash with disabled wrapper",
			shell:         "bash",
			functionAlias: "false",
			wantWrapper:   false,
			wantFunction:  "",
		},
		{
			name:          "fish with default function name",
			shell:         "fish",
			functionAlias: "true",
			wantWrapper:   true,
			wantFunction:  "function bookmark",
		},
		{
			name:          "fish with custom function name",
			shell:         "fish",
			functionAlias: "bm",
			wantWrapper:   true,
			wantFunction:  "function bm",
		},
		{
			name:          "zsh with default function name",
			shell:         "zsh",
			functionAlias: "true",
			wantWrapper:   true,
			wantFunction:  "bookmark() {",
		},
		{
			name:          "nushell with default function name",
			shell:         "nu",
			functionAlias: "true",
			wantWrapper:   true,
			wantFunction:  "def --wrapped bookmark",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "bookmarks.sh")

			m := NewManager(filePath, tt.shell, "cd", "nvim", tt.functionAlias, "bm")

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

			// Check for function wrapper presence
			hasWrapper := strings.Contains(contentStr, "Function wrapper to auto-source")
			if hasWrapper != tt.wantWrapper {
				t.Errorf("wrapper presence = %v, want %v", hasWrapper, tt.wantWrapper)
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
		})
	}
}

func TestManager_FunctionWrapperSourcesCorrectFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bookmarks.sh")

	m := NewManager(filePath, "bash", "cd", "nvim", "true", "bm")

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

	// Verify the wrapper sources the correct file path
	expectedSource := "source " + filePath
	if !strings.Contains(contentStr, expectedSource) {
		t.Errorf("expected source command %q not found in:\n%s", expectedSource, contentStr)
	}

	// Verify the wrapper uses command to avoid recursion
	if !strings.Contains(contentStr, "command bookmark") {
		t.Error("expected 'command bookmark' to avoid recursion")
	}

	// Verify conditional execution with &&
	if !strings.Contains(contentStr, "&&") {
		t.Error("expected && to only source on successful command execution")
	}
}

func TestManager_FunctionWrapperEmptyAlias(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bookmarks.sh")

	m := NewManager(filePath, "bash", "cd", "nvim", "", "bm")

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

	// Empty string should disable wrapper
	if strings.Contains(contentStr, "Function wrapper") {
		t.Error("wrapper should not be generated when functionAlias is empty")
	}
}
