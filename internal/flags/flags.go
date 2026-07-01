package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Set registers a flag on cmd inferred from the type of ptr.
// Supported types: *string (default ""), *bool (default false), *int (default 0).
// group annotates the flag for help grouping; pass "" to skip.
// If you want to define flags differently you can simply call cmd.flags directly
func Set(cmd *cobra.Command, ptr any, name, shorthand, usage, group string) {
	register(cmd.Flags(), ptr, name, shorthand, usage)
	if group != "" {
		_ = cmd.Flags().SetAnnotation(name, "group", []string{group})
	}
}

// SetPersistent registers a persistent flag on cmd (inherited by subcommands).
func SetPersistent(cmd *cobra.Command, ptr any, name, shorthand, usage, group string) {
	register(cmd.PersistentFlags(), ptr, name, shorthand, usage)
	if group != "" {
		_ = cmd.PersistentFlags().SetAnnotation(name, "group", []string{group})
	}
}

// register registers a flag based on the pointer type passed.
func register(fs *pflag.FlagSet, ptr any, name, shorthand, usage string) {
	switch p := ptr.(type) {
	case *string:
		fs.StringVarP(p, name, shorthand, "", usage)
	case *bool:
		fs.BoolVarP(p, name, shorthand, false, usage)
	case *int:
		fs.IntVarP(p, name, shorthand, 0, usage)
	}
}
