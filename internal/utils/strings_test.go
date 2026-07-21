package utils

import "testing"

func TestGenerateAlias(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		separator  string
		lowercase  bool
		partLength int
		want       string
	}{
		{
			name:       "simple path - first letter",
			path:       "/home/user/my-cool-project",
			separator:  "",
			lowercase:  true,
			partLength: 1,
			want:       "mcp",
		},
		{
			name:       "simple path - two letters",
			path:       "/home/user/my-cool-project",
			separator:  "",
			lowercase:  true,
			partLength: 2,
			want:       "mycopr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateAlias(tt.path, tt.separator, tt.lowercase, tt.partLength)
			if got != tt.want {
				t.Errorf("GenerateAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEscapeAliasValue(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want string
	}{
		{
			name: "no quotes",
			val:  "echo hi",
			want: "echo hi",
		},
		{
			name: "with single quotes",
			val:  "echo 'hi'",
			want: "echo 'hi'",
		},
		{
			name: "with double quotes",
			val:  `echo "hi"`,
			want: "echo 'hi'",
		},
		{
			name: "with mixed quotes",
			val:  `echo "hello" and 'world'`,
			want: "echo 'hello' and 'world'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapeAliasValue(tt.val)
			if got != tt.want {
				t.Errorf("EscapeAliasValue() = %q, want %q", got, tt.want)
			}
		})
	}
}
