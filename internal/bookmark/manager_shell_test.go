package bookmark

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aliases/internal/domain"
)

func TestNushellAliasFormat(t *testing.T) {
	tempDir := t.TempDir()
	bookmarkFile := filepath.Join(tempDir, "bookmarks.nu")

	mgr := NewManager(bookmarkFile, "nu", "cd", "vim", "true", "bm")
	
	t.Logf("Manager shell: %q", mgr.shell)

	bm := domain.Bookmark{
		Alias:     "test",
		Path:      "/tmp/test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := mgr.Add(bm)
	if err != nil {
		t.Fatalf("failed to add bookmark: %v", err)
	}

	content, err := os.ReadFile(bookmarkFile)
	if err != nil {
		t.Fatalf("failed to read bookmark file: %v", err)
	}

	contentStr := string(content)
	
	// Check for nushell format with spaces around =
	if !contains(contentStr, "alias test = ") {
		t.Errorf("Expected nushell format 'alias test = ', got:\n%s", contentStr)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
