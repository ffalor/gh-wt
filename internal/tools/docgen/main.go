package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ffalor/gh-wt/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	out := flag.String("out", "./docs/cli", "output directory")
	_ = flag.String("format", "markdown", "output format (markdown only)")
	frontmatter := flag.Bool("frontmatter", false, "include frontmatter")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}

	root := cmd.Root()
	root.DisableAutoGenTag = true

	if err := genMarkdownTree(root, *out, *frontmatter); err != nil {
		log.Fatal(err)
	}
}

// genMarkdownTree generates markdown docs for the command tree.
func genMarkdownTree(cmd *cobra.Command, outDir string, frontmatter bool) error {
	preamble := func(filename string) string {
		if !frontmatter {
			return ""
		}
		base := filepath.Base(filename)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		title := strings.ReplaceAll(name, "_", " ")
		return fmt.Sprintf("---\ntitle: %q\nlayout: '../../../layouts/CLILayout.astro'\nslug: %q\ndescription: \"CLI reference for %s\"\n---\n\n", title, name, title)
	}

	seen := make(map[string]bool)
	return walkCMD(cmd, func(cc *cobra.Command) error {
		// Generate filename
		name := cc.CommandPath()
		name = strings.ReplaceAll(name, " ", "_")
		name = strings.ToLower(name) + ".mdx"
		filename := filepath.Join(outDir, name)

		// Skip if already generated (e.g., help command)
		if seen[filename] {
			return nil
		}
		seen[filename] = true

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
			return err
		}

		// Generate content
		buf := bytes.NewBufferString(preamble(filename))
		if err := genMarkdown(buf, cc); err != nil {
			return err
		}

		return os.WriteFile(filename, buf.Bytes(), 0o644)
	})
}

// walkCMD walks the command tree.
func walkCMD(cmd *cobra.Command, fn func(*cobra.Command) error) error {
	// Generate docs for all commands except help
	if cmd.Name() != "help" {
		if err := fn(cmd); err != nil {
			return err
		}
	}
	for _, c := range cmd.Commands() {
		if err := walkCMD(c, fn); err != nil {
			return err
		}
	}
	return nil
}

// genMarkdown generates markdown for a single command.
func genMarkdown(buf *bytes.Buffer, cmd *cobra.Command) error {
	title := cmd.CommandPath()
	buf.WriteString("## " + title + "\n\n")

	if cmd.Short != "" {
		buf.WriteString(cmd.Short + "\n\n")
	}

	// Synopsis with Long description and usage
	buf.WriteString("### Synopsis\n\n")
	if cmd.Long != "" {
		buf.WriteString(cmd.Long + "\n\n")
	}
	synopsis := cmd.CommandPath()
	if cmd.Parent() != nil {
		synopsis = cmd.UseLine()
	}
	fmt.Fprintf(buf, "```\n%s\n```\n\n", synopsis)

	// Examples
	if cmd.Example != "" {
		buf.WriteString("### Examples\n\n")
		buf.WriteString("```\n")
		buf.WriteString(cmd.Example)
		buf.WriteString("```\n\n")
	}

	// Options - separate tables for local and inherited options
	localFlags := cmd.NonInheritedFlags()
	inheritedFlags := cmd.InheritedFlags()

	// Output local flags table
	if localFlags.HasFlags() {
		buf.WriteString("### Options\n\n")

		type flagInfo struct {
			names string
			usage string
		}

		var flags []flagInfo
		localFlags.VisitAll(func(f *pflag.Flag) {
			flags = append(flags, flagInfo{
				names: formatFlagNames(f),
				usage: f.Usage,
			})
		})

		sort.Slice(flags, func(i, j int) bool {
			return flags[i].names < flags[j].names
		})

		buf.WriteString("| Flag | Description |\n")
		buf.WriteString("|------|-------------|\n")
		for _, f := range flags {
			desc := strings.ReplaceAll(f.usage, "|", "\\|")
			fmt.Fprintf(buf, "| `%s` | %s |\n", f.names, desc)
		}
		buf.WriteString("\n")
	}

	// Output inherited flags table
	if inheritedFlags.HasFlags() {
		buf.WriteString("### Options inherited from parent commands\n\n")

		type flagInfo struct {
			names string
			usage string
		}

		var flags []flagInfo
		inheritedFlags.VisitAll(func(f *pflag.Flag) {
			flags = append(flags, flagInfo{
				names: formatFlagNames(f),
				usage: f.Usage,
			})
		})

		sort.Slice(flags, func(i, j int) bool {
			return flags[i].names < flags[j].names
		})

		buf.WriteString("| Flag | Description |\n")
		buf.WriteString("|------|-------------|\n")
		for _, f := range flags {
			desc := strings.ReplaceAll(f.usage, "|", "\\|")
			fmt.Fprintf(buf, "| `%s` | %s |\n", f.names, desc)
		}
		buf.WriteString("\n")
	}

	// Additional help topics
	if cmd.HasAvailableSubCommands() {
		buf.WriteString("### SEE ALSO\n\n")

		var links []string
		for _, c := range cmd.Commands() {
			if c.Name() == "help" {
				continue
			}
			link := strings.ReplaceAll(c.CommandPath(), " ", "_")
			linkLower := strings.ToLower(link)
			links = append(links, fmt.Sprintf("* [%s](/docs/cli/%s)\t - %s", link, linkLower, c.Short))
		}
		buf.WriteString(strings.Join(links, "\n\n"))
		buf.WriteString("\n")
	}

	return nil
}

// formatFlagNames formats flag names for display.
func formatFlagNames(f *pflag.Flag) string {
	names := ""
	if f.Shorthand != "" {
		names = "-" + f.Shorthand + ", --" + f.Name
	} else {
		names = "--" + f.Name
	}
	if f.Value.Type() != "bool" && f.DefValue != "" {
		names += " " + strings.ReplaceAll(f.DefValue, "\n", "")
	}
	return names
}
