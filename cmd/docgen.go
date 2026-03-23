//go:build docgen

// docgen is a build-tag-gated subcommand that generates the CLI reference page
// for the giftwrap Hugo documentation site.
//
// Usage: go run -tags docgen . docgen
// Output: _doc/content/docs/cli.md (default)
package cmd

import (
	_ "embed"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

//go:embed docgen_remarks.md
var docgenRemarks []byte

var docgenOut *string

var docgenCmd = &cobra.Command{
	Use:   "docgen",
	Short: "Generate the CLI reference page for the Hugo docs site",
	Long:  `docgen writes a manpage-style CLI reference to the Hugo docs tree. Build with -tags docgen.`,
	// Suppress the root PersistentPreRunE — docgen doesn't need the logger or runner setup.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	Run:               runDocgen,
}

func init() {
	rootCmd.AddCommand(docgenCmd)
	docgenOut = docgenCmd.Flags().String("output", "_doc/content/docs/cli.md", "Output path for the generated CLI reference")
}

func runDocgen(cmd *cobra.Command, args []string) {
	root := RootCommand()

	var buf bytes.Buffer

	buf.WriteString("---\ntitle: \"giftwrap(1)\"\ndraft: false\n---\n\n")

	// Synopsis
	buf.WriteString("## Synopsis\n\n")
	buf.WriteString("**giftwrap** — cross-compile Go applications and package releases.\n\n")
	if long := normalizeLong(root.Long); long != "" {
		buf.WriteString(long)
		buf.WriteString("\n\n")
	}

	// Usage
	buf.WriteString("## Usage\n\n")
	buf.WriteString("    giftwrap [--wrapfile path] [--log-level level] <command> [options]\n\n")
	visible := docgenVisibleCmds(root)
	buf.WriteString("Available commands:\n\n")
	for _, c := range visible {
		fmt.Fprintf(&buf, "    %-12s %s\n", c.Name(), c.Short)
	}
	buf.WriteString("\n")

	// Global Options
	buf.WriteString("## Global Options\n\n")
	helpEntry := docgenFlagEntry{left: "-h, --help", right: "Print help"}
	buf.WriteString(docgenRenderEntries(append([]docgenFlagEntry{helpEntry}, docgenFlagEntries(root.PersistentFlags())...)))
	buf.WriteString("\n")

	// Commands
	buf.WriteString("## Commands\n\n")
	for i, c := range visible {
		fmt.Fprintf(&buf, "### %s — %s\n\n", c.Name(), c.Short)
		if long := normalizeLong(c.Long); long != "" {
			buf.WriteString(long)
			buf.WriteString("\n\n")
		}
		if flags := docgenRenderEntries(docgenFlagEntries(c.Flags())); flags != "" {
			buf.WriteString("**Options:**\n\n")
			buf.WriteString(flags)
			buf.WriteString("\n")
		}
		if i < len(visible)-1 {
			buf.WriteString("---\n\n")
		}
	}

	// Remarks
	if remarks := strings.TrimSpace(string(docgenRemarks)); remarks != "" {
		buf.WriteString("\n## Remarks\n\n")
		buf.WriteString(remarks)
		buf.WriteString("\n")
	}

	out := *docgenOut
	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(out, buf.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", out, err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s\n", out)
}

// cobraBuiltins are commands auto-added by Cobra that we exclude from the manpage.
var cobraBuiltins = map[string]bool{"completion": true, "help": true, "docgen": true}

// docgenVisibleCmds returns user-defined, non-hidden subcommands sorted by name.
func docgenVisibleCmds(root *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	for _, c := range root.Commands() {
		if !c.Hidden && !cobraBuiltins[c.Name()] {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out
}

// normalizeLong trims leading whitespace/tabs from each line of a Long description.
func normalizeLong(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimLeft(l, "\t ")
	}
	return strings.Join(lines, "\n")
}

type docgenFlagEntry struct {
	left  string
	right string
}

func docgenFlagEntries(fs *pflag.FlagSet) []docgenFlagEntry {
	var entries []docgenFlagEntry
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Name == "help" {
			return
		}
		var left string
		if f.Shorthand != "" {
			left = fmt.Sprintf("-%s, --%s", f.Shorthand, f.Name)
		} else {
			left = fmt.Sprintf("    --%s", f.Name)
		}
		if f.Value.Type() != "bool" {
			left += " " + f.Value.Type()
		}
		right := f.Usage
		if f.DefValue != "" && f.DefValue != "false" {
			right += fmt.Sprintf(" (default: %s)", f.DefValue)
		}
		entries = append(entries, docgenFlagEntry{left, right})
	})
	return entries
}

func docgenRenderEntries(entries []docgenFlagEntry) string {
	if len(entries) == 0 {
		return ""
	}
	maxLeft := 0
	for _, e := range entries {
		if len(e.left) > maxLeft {
			maxLeft = len(e.left)
		}
	}
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "    %-*s   %s\n", maxLeft, e.left, e.right)
	}
	return sb.String()
}
