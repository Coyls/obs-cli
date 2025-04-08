package push

import (
	"fmt"

	"github.com/coyls/obs-cli/internal/config"
	"github.com/coyls/obs-cli/internal/git"
	"github.com/coyls/obs-cli/internal/logger"
	"github.com/spf13/cobra"
)

var (
	force bool
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes to GitHub",
	Long: `The push command synchronizes your Obsidian vault with the remote GitHub repository.
It performs a commit and pushes the changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executePush()
	},
}

func executePush() error {
	logger.PrintHeader("Push Obsidian to GitHub")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	gitClient := git.New(cfg.Config.Root)

	logger.Info("Checking for changes...")
	hasChanges, err := gitClient.HasChanges()
	if err != nil {
		logger.Error("%s", err.Error())
		return err
	}

	if !hasChanges && !force {
		logger.Info("No changes to add")
		return nil
	}

	// Add changes
	logger.Info("Adding changes...")
	if err := gitClient.AddAll(); err != nil {
		logger.Error("%s", err.Error())
		return err
	}
	logger.Success("Changes added")

	// Create commit
	logger.Info("Creating commit...")
	if err := gitClient.Commit(); err != nil {
		logger.Info("No changes to commit")
		return nil
	}
	logger.Success("Commit created")

	// Push to GitHub
	logger.Info("Pushing to GitHub...")
	if err := gitClient.Push(); err != nil {
		logger.Error("Unable to connect to GitHub")
		logger.Error("Check your internet connection and try again")
		return err
	}
	logger.Success("Push successful!")

	logger.Success("Synchronization completed!")
	return nil
}

func init() {
	pushCmd.Flags().BoolVarP(&force, "force", "f", false, "Force push even without changes")
}

// GetCommand returns the push command for root command integration
func GetCommand() *cobra.Command {
	return pushCmd
}
