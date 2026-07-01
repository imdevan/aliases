package flags_test

import (
	"testing"

	"github.com/spf13/cobra"

	"bookmark/internal/flags"
)

func newCmd() *cobra.Command {
	return &cobra.Command{Use: "test"}
}

func TestSet_string(t *testing.T) {
	cmd := newCmd()
	var val string
	flags.Set(cmd, &val, "name", "n", "a name", "")

	f := cmd.Flags().Lookup("name")
	if f == nil {
		t.Fatal("flag not registered")
	}
	if f.DefValue != "" {
		t.Errorf("default = %q, want \"\"", f.DefValue)
	}
	if f.Usage != "a name" {
		t.Errorf("usage = %q, want \"a name\"", f.Usage)
	}
	if f.Shorthand != "n" {
		t.Errorf("shorthand = %q, want \"n\"", f.Shorthand)
	}
}

func TestSet_bool(t *testing.T) {
	cmd := newCmd()
	var val bool
	flags.Set(cmd, &val, "verbose", "v", "enable verbose output", "")

	f := cmd.Flags().Lookup("verbose")
	if f == nil {
		t.Fatal("flag not registered")
	}
	if f.DefValue != "false" {
		t.Errorf("default = %q, want \"false\"", f.DefValue)
	}
	if f.Usage != "enable verbose output" {
		t.Errorf("usage = %q, want \"enable verbose output\"", f.Usage)
	}
}

func TestSet_int(t *testing.T) {
	cmd := newCmd()
	var val int
	flags.Set(cmd, &val, "count", "c", "number of items", "")

	f := cmd.Flags().Lookup("count")
	if f == nil {
		t.Fatal("flag not registered")
	}
	if f.DefValue != "0" {
		t.Errorf("default = %q, want \"0\"", f.DefValue)
	}
}

func TestSet_groupAnnotation(t *testing.T) {
	cmd := newCmd()
	var val string
	flags.Set(cmd, &val, "config", "c", "config file", "config")

	annotations := cmd.Flags().Lookup("config").Annotations
	group, ok := annotations["group"]
	if !ok {
		t.Fatal("group annotation not set")
	}
	if len(group) != 1 || group[0] != "config" {
		t.Errorf("group = %v, want [\"config\"]", group)
	}
}

func TestSet_noGroupAnnotation(t *testing.T) {
	cmd := newCmd()
	var val string
	flags.Set(cmd, &val, "name", "n", "a name", "")

	f := cmd.Flags().Lookup("name")
	if _, ok := f.Annotations["group"]; ok {
		t.Error("expected no group annotation when group is empty")
	}
}

func TestSetPersistent(t *testing.T) {
	cmd := newCmd()
	var val string
	flags.SetPersistent(cmd, &val, "config", "c", "config file", "config")

	if cmd.Flags().Lookup("config") != nil {
		t.Error("flag should not appear in Flags(), only PersistentFlags()")
	}
	f := cmd.PersistentFlags().Lookup("config")
	if f == nil {
		t.Fatal("flag not registered in PersistentFlags")
	}
	if f.Usage != "config file" {
		t.Errorf("usage = %q, want \"config file\"", f.Usage)
	}
	group := f.Annotations["group"]
	if len(group) != 1 || group[0] != "config" {
		t.Errorf("group = %v, want [\"config\"]", group)
	}
}
