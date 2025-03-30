package initcmd

import (
	"os"
	"path/filepath"

	"github.com/coyls/obs-cli/internal/config"
	"github.com/coyls/obs-cli/internal/logger"
	"github.com/spf13/cobra"
)

var (
	force bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize or update the configuration",
	Long: `The init command creates or updates the configuration file.
By default, it will create a new configuration file if it doesn't exist.
Use the --force flag to overwrite an existing configuration file.

The configuration file is located at ~/.config/obs-cli/config.yaml
You can edit this file manually to customize the settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeInit()
	},
}

func executeInit() error {
	logger.PrintHeader("Initialize configuration")

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Failed to get home directory: %s", err.Error())
		return err
	}

	// Define config file path
	configPath := filepath.Join(homeDir, ".config", "obs-cli", "config.yaml")

	// Check if config file exists
	if _, err := os.Stat(configPath); err == nil && !force {
		logger.Info("Configuration file already exists at: %s", configPath)
		logger.Info("To modify the configuration, you can:")
		logger.Info("1. Edit the file manually: %s", configPath)
		logger.Info("2. Use --force to overwrite the existing configuration")
		return nil
	}

	// Create default config
	cfg := config.DefaultConfig()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		logger.Error("Failed to create config directory: %s", err.Error())
		return err
	}

	// Save config
	if err := cfg.Save(configPath); err != nil {
		logger.Error("Failed to save configuration: %s", err.Error())
		return err
	}

	logger.Success("Configuration file created at: %s", configPath)
	logger.Info("\nCurrent configuration:")
	logger.Info("Vault path: %s", cfg.VaultPath)
	logger.Info("Default move path: %s", cfg.DefaultMvPath)
	logger.Info("\nTo modify the configuration, edit the file manually or use --force to overwrite it.")
	return nil
}

func init() {
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite of existing configuration")
}

// GetCommand returns the init command for root command integration
func GetCommand() *cobra.Command {
	return initCmd
}
