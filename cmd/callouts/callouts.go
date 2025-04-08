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

func executeCallouts(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	vaultConfig, exists := cfg.GetVaultConfig(cfg.Config.DefaultVault)
	if !exists {
		logger.Error("Default vault configuration not found")
		return fmt.Errorf("default vault configuration not found")
	}

	calloutsPath := filepath.Join(cfg.Config.Root, vaultConfig.VaultPath, ".obsidian", "snippets", "snippet.css")

	if _, err := os.Stat(calloutsPath); os.IsNotExist(err) {
		snippetsDir := filepath.Dir(calloutsPath)
		if err := os.MkdirAll(snippetsDir, 0755); err != nil {
			return fmt.Errorf("failed to create snippets directory: %w", err)
		}

		if err := os.WriteFile(calloutsPath, []byte("/* Add your callout styles here */\n"), 0644); err != nil {
			return fmt.Errorf("failed to create callouts file: %w", err)
		}

		logger.Info("Created new callouts file at: %s", calloutsPath)
	}

	editor := cfg.Config.DefaultEditor
	if editor == "" {
		if envEditor := os.Getenv("EDITOR"); envEditor != "" {
			editor = envEditor
		} else {
			editor = "nano"
		}
	}

	editorCmd := exec.Command(editor, calloutsPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	logger.Success("Callouts file opened successfully!")
	return nil
}
