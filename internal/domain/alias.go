package domain

// Alias represents a command alias with metadata.
type Alias struct {
	Name        string `toml:"name"`
	Value       string `toml:"value"`
	Description string `toml:"description,omitempty"`
	SourceFile  string `toml:"source_file,omitempty"`
	Global      bool   `toml:"global,omitempty"`
	Suffix      bool   `toml:"suffix,omitempty"`
}
