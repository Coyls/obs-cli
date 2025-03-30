package callouts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coyls/obs-cli/internal/config"
	"github.com/coyls/obs-cli/internal/logger"
	"github.com/spf13/cobra"
)

// GetCommand returns the callouts command
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "callouts",
		Short: "Edit Obsidian callouts",
		Long: `The callouts command opens your Obsidian callouts configuration file in your default editor.
This allows you to customize the appearance and behavior of callouts in your vault.`,
		RunE: executeCallouts,
	}

	return cmd
}

// executeCallouts handles the callouts command execution
func executeCallouts(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Build path to callouts file
	calloutsPath := filepath.Join(cfg.ObsidianPath, "snippets", "snippet.css")

	// Check if callouts file exists
	if _, err := os.Stat(calloutsPath); os.IsNotExist(err) {
		// Create snippets directory if it doesn't exist
		snippetsDir := filepath.Dir(calloutsPath)
		if err := os.MkdirAll(snippetsDir, 0755); err != nil {
			return fmt.Errorf("failed to create snippets directory: %w", err)
		}

		// Create empty callouts file
		if err := os.WriteFile(calloutsPath, []byte("/* Add your callout styles here */\n"), 0644); err != nil {
			return fmt.Errorf("failed to create callouts file: %w", err)
		}

		logger.Info("Created new callouts file at: %s", calloutsPath)
	}

	// Open file in default editor
	editorCmd := exec.Command(cfg.DefaultEditor, calloutsPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	logger.Success("Callouts file opened successfully!")
	return nil
}
