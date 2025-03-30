package main

import (
	"os"

	"github.com/coyls/obs-cli/cmd/callouts"
	"github.com/coyls/obs-cli/cmd/cp"
	"github.com/coyls/obs-cli/cmd/initcmd"
	"github.com/coyls/obs-cli/cmd/mv"
	"github.com/coyls/obs-cli/cmd/pull"
	"github.com/coyls/obs-cli/cmd/push"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "obs-cli",
	Short: "CLI to manage your Obsidian vault",
	Long: `obs-cli is a command-line interface for managing your Obsidian vault.
It provides commands to synchronize your vault with Git and organize your notes.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	// Add commands
	rootCmd.AddCommand(initcmd.GetCommand())
	rootCmd.AddCommand(push.GetCommand())
	rootCmd.AddCommand(pull.GetCommand())
	rootCmd.AddCommand(mv.GetCommand())
	rootCmd.AddCommand(cp.GetCommand())
	rootCmd.AddCommand(callouts.GetCommand())
	Execute()
}
