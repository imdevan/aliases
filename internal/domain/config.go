package domain

import (
	"os"
	"path/filepath"

	shelladapter "bookmark/internal/adapters/shell"
)

// Config describes the resolved configuration.
type Config struct {
	Editor               string `toml:"editor"`
	Primary              string `toml:"primary"`
	Secondary            string `toml:"secondary"`
	Headings             string `toml:"headings"`
	Text                 string `toml:"text"`
	TextHighlight        string `toml:"text_highlight"`
	DescriptionHighlight string `toml:"description_highlight"`
	Tags                 string `toml:"tags"`
	Flags                string `toml:"flags"`
	Muted                string `toml:"muted"`
	Accent               string `toml:"accent"`
	Border               string `toml:"border"`
	Error                string `toml:"error"`
	InteractiveDefault   bool   `toml:"interactive_default"`
	ListSpacing          string `toml:"list_spacing"`
	
	// Bookmark settings
	BookmarkLocation       string `toml:"bookmark_location"`
	NavigationTool         string `toml:"navigation_tool"`
	Shell                  string `toml:"shell"`
	AutoAliasSeparator     string `toml:"auto_alias_separator"`
	AutoAliasLowercase     bool   `toml:"auto_alias_lowercase"`
	DefaultAliasPartLength int    `toml:"default_alias_part_length"`
	HomeIcon               string `toml:"home_icon"`
	DefaultSortBy          string `toml:"default_sort_by"`
	FunctionAlias          string `toml:"function_alias"`
	InteractiveAlias       string `toml:"interactive_alias"`
}

// DefaultConfig returns the default configuration values.
func DefaultConfig() Config {
	home, _ := os.UserHomeDir()
	bookmarkLocation := filepath.Join(home, ".bookmarks")
	detectedShell := shelladapter.DetectShell()
	
	return Config{
		Editor:               "nvim",
		Headings:             "15",
		Primary:              "02",
		Secondary:            "06",
		Text:                 "07",
		TextHighlight:        "06",
		DescriptionHighlight: "05",
		Tags:                 "13",
		Flags:                "12",
		Muted:                "08",
		Accent:               "13",
		Border:               "08",
		Error:                "01",
		InteractiveDefault:   false,
		ListSpacing:          "space",
		BookmarkLocation:     bookmarkLocation,
		NavigationTool:       "cd",
		Shell:                  detectedShell,
		AutoAliasSeparator:     "",
		AutoAliasLowercase:     true,
		DefaultAliasPartLength: 1, // Take 1 character from each part by default
		HomeIcon:               "~",
		DefaultSortBy:        "newest",
		FunctionAlias:        "true",
		InteractiveAlias:     "bm",
	}
}

// GetBookmarkFileName returns the appropriate bookmark filename for the shell.
func GetBookmarkFileName(shell string) string {
	switch shell {
	case "fish":
		return "bookmarks.fish"
	case "nu", "nushell":
		return "bookmarks.nu"
	default: // bash, zsh, sh
		return "bookmarks.sh"
	}
}

// BookmarkFile returns the full path to the bookmark file based on shell.
func (c Config) BookmarkFile() string {
	return filepath.Join(c.BookmarkLocation, GetBookmarkFileName(c.Shell))
}

func xdgHome(envKey, fallbackSuffix string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, fallbackSuffix)
}
