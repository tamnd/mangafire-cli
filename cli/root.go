// Package cli builds the mangafire command tree on top of the mangafire library.
package cli

import (
	"github.com/spf13/cobra"
)

// Build metadata, set via -ldflags at release time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// Root builds the root command and its subtree.
func Root() *cobra.Command {
	root := &cobra.Command{
		Use:   "mangafire",
		Short: "A command line for mangafire.",
		Long: `A command line for mangafire.

This is a fresh scaffold. Add your commands here on top of the mangafire
library package, then wire them into Root with root.AddCommand.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(newVersionCmd())
	// TODO: root.AddCommand(newGetCmd()), etc.
	return root
}
