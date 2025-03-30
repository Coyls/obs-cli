package pull

import (
	"strings"

	"github.com/coyls/obs-cli/internal/config"
	"github.com/coyls/obs-cli/internal/git"
	"github.com/coyls/obs-cli/internal/logger"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from GitHub",
	Long: `The pull command synchronizes your Obsidian vault with the remote GitHub repository.
It fetches and applies the latest changes locally.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executePull()
	},
}

func executePull() error {
	logger.PrintHeader("Pull Obsidian from GitHub")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Initialize Git client
	gitClient := git.New(cfg.VaultPath)

	// Check current branch
	logger.Info("Checking current branch...")
	currentBranch, err := gitClient.GetCurrentBranch()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if currentBranch != "main" {
		logger.Error("You are not on the main branch (current branch: %s)", currentBranch)
		return nil
	}

	// Check remote changes
	logger.Info("Checking remote changes...")
	if err := gitClient.Fetch(); err != nil {
		logger.Error("Error while checking for changes")
		return err
	}
	logger.Success("Check completed")

	// Pull changes
	logger.Info("Fetching changes...")
	if err := gitClient.Pull(); err != nil {
		// Handle conflicts
		if strings.Contains(err.Error(), "conflict") {
			logger.Error("Conflicts detected!")
			logger.Info("Conflicting files:")

			// List conflicting files
			conflicts, err := gitClient.GetConflicts()
			if err != nil {
				logger.Error("Unable to list conflicts")
				return err
			}

			for _, conflict := range conflicts {
				logger.Info("  - %s", conflict)
			}

			logger.Info("\nResolve conflicts manually and commit changes")
		} else {
			logger.Error(err.Error())
		}
		return err
	}

	logger.Success("Pull successful!")
	logger.Success("Synchronization completed!")
	return nil
}

// GetCommand returns the pull command for root command integration
func GetCommand() *cobra.Command {
	return pullCmd
}
