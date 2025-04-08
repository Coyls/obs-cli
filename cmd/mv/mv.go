package mv

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coyls/obs-cli/internal/config"
	"github.com/coyls/obs-cli/internal/logger"
	"github.com/spf13/cobra"
)

var (
	destination string
)

var mvCmd = &cobra.Command{
	Use:   "mv [source_file]",
	Short: "Move a file to the Obsidian vault",
	Long: `The mv command moves a file from anywhere on your system to your Obsidian vault.
If no destination is specified, the file will be moved to the default directory defined in the configuration.

Example:
  obs-cli mv ~/Downloads/image.png -d Assets/new`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return executeMove(args[0])
	},
}

func executeMove(source string) error {
	logger.PrintHeader("Move file to Obsidian vault")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if _, err := os.Stat(source); os.IsNotExist(err) {
		logger.Error("Source file not found: %s", source)
		return fmt.Errorf("source file not found: %s", source)
	}

	vaultConfig, exists := cfg.GetVaultConfig(cfg.Config.DefaultVault)
	if !exists {
		logger.Error("Default vault configuration not found")
		return fmt.Errorf("default vault configuration not found")
	}

	if destination == "" {
		if vaultConfig.Commands.Mv.DefaultTargetPath == "" {
			logger.Error("No destination specified and no default path configured")
			return fmt.Errorf("no destination specified and no default path configured")
		}
		destination = vaultConfig.Commands.Mv.DefaultTargetPath
		logger.Info("Using default destination: %s", destination)
	}

	destPath := filepath.Join(cfg.Config.Root, vaultConfig.VaultPath, destination)
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		logger.Info("Creating destination directory: %s", destPath)
		if err := os.MkdirAll(destPath, 0755); err != nil {
			logger.Error("Failed to create destination directory: %s", err.Error())
			return fmt.Errorf("failed to create destination directory: %w", err)
		}
	}

	filename := filepath.Base(source)
	finalDest := filepath.Join(destPath, filename)

	if _, err := os.Stat(finalDest); err == nil {
		logger.Error("File already exists in destination: %s", finalDest)
		return fmt.Errorf("file already exists in destination: %s", finalDest)
	}

	logger.Info("Moving file from %s to %s...", source, finalDest)
	if err := os.Rename(source, finalDest); err != nil {
		logger.Error("Failed to move file: %s", err.Error())
		return fmt.Errorf("failed to move file: %w", err)
	}

	logger.Success("File moved successfully!")
	return nil
}

func init() {
	mvCmd.Flags().StringVarP(&destination, "destination", "d", "", "Destination directory in the vault (optional)")
}

func GetCommand() *cobra.Command {
	return mvCmd
}
