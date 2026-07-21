package editor

import (
	"strings"
	"testing"
)

func TestBuildCmd(t *testing.T) {
	tests := []struct {
		name      string
		editor    string
		path      string
		line      int
		wantCmd   string
		wantArgs  []string
	}{
		{
			name:     "vim with line number",
			editor:   "vim",
			path:     "/tmp/aliases.zsh",
			line:     15,
			wantCmd:  "vim",
			wantArgs: []string{"+15", "/tmp/aliases.zsh"},
		},
		{
			name:     "code with line number",
			editor:   "code",
			path:     "/tmp/aliases.zsh",
			line:     42,
			wantCmd:  "code",
			wantArgs: []string{"-g", "/tmp/aliases.zsh:42"},
		},
		{
			name:     "nano with line number",
			editor:   "nano",
			path:     "/tmp/aliases.zsh",
			line:     10,
			wantCmd:  "nano",
			wantArgs: []string{"+10", "/tmp/aliases.zsh"},
		},
		{
			name:     "editor without line number",
			editor:   "vim",
			path:     "/tmp/aliases.zsh",
			line:     0,
			wantCmd:  "vim",
			wantArgs: []string{"/tmp/aliases.zsh"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := New(tt.editor)
			cmd, err := adapter.BuildCmd(tt.path, tt.line)
			if err != nil {
				t.Fatalf("BuildCmd error = %v", err)
			}
			if cmd.Path != tt.wantCmd && !strings.HasSuffix(cmd.Path, "/"+tt.wantCmd) {
				t.Errorf("cmd.Path = %q, want %q", cmd.Path, tt.wantCmd)
			}
			if len(cmd.Args)-1 != len(tt.wantArgs) {
				t.Fatalf("len(cmd.Args) = %d, want %d (args: %v)", len(cmd.Args)-1, len(tt.wantArgs), cmd.Args)
			}
			for i, arg := range tt.wantArgs {
				if cmd.Args[i+1] != arg {
					t.Errorf("arg[%d] = %q, want %q", i, cmd.Args[i+1], arg)
				}
			}
		})
	}
}
