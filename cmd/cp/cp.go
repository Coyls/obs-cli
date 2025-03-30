package cp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/coyls/obs-cli/internal/config"
	"github.com/coyls/obs-cli/internal/logger"
	"github.com/spf13/cobra"
)

var (
	destination string
)

var cpCmd = &cobra.Command{
	Use:   "cp [source_file]",
	Short: "Copy a file to the Obsidian vault",
	Long: `The cp command copies a file from anywhere on your system to your Obsidian vault.
If no destination is specified, the file will be copied to the default directory defined in the configuration.

Example:
  obs-cli cp ~/Downloads/image.png -d Assets/new`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeCopy(args[0])
	},
}

func executeCopy(source string) error {
	logger.PrintHeader("Copy file to Obsidian vault")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("%s", err.Error())
		return err
	}

	// Check if source file exists
	if _, err := os.Stat(source); os.IsNotExist(err) {
		logger.Error("Source file not found: %s", source)
		return fmt.Errorf("source file not found: %s", source)
	}

	// Use default destination if none specified
	if destination == "" {
		if cfg.DefaultCpPath == "" {
			logger.Error("No destination specified and no default path configured")
			return fmt.Errorf("no destination specified and no default path configured")
		}
		destination = cfg.DefaultCpPath
		logger.Info("Using default destination: %s", destination)
	}

	// Build destination path in the vault
	destPath := filepath.Join(cfg.VaultPath, destination)
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		logger.Info("Creating destination directory: %s", destPath)
		if err := os.MkdirAll(destPath, 0755); err != nil {
			logger.Error("Failed to create destination directory: %s", err.Error())
			return fmt.Errorf("failed to create destination directory: %w", err)
		}
	}

	// Get the filename from source
	filename := filepath.Base(source)
	finalDest := filepath.Join(destPath, filename)

	// Check if destination file already exists
	if _, err := os.Stat(finalDest); err == nil {
		logger.Error("File already exists in destination: %s", finalDest)
		return fmt.Errorf("file already exists in destination: %s", finalDest)
	}

	// Open source file
	srcFile, err := os.Open(source)
	if err != nil {
		logger.Error("Failed to open source file: %s", err.Error())
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(finalDest)
	if err != nil {
		logger.Error("Failed to create destination file: %s", err.Error())
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy the file
	logger.Info("Copying file from %s to %s...", source, finalDest)
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		logger.Error("Failed to copy file: %s", err.Error())
		return fmt.Errorf("failed to copy file: %w", err)
	}

	logger.Success("File copied successfully!")
	return nil
}

func init() {
	cpCmd.Flags().StringVarP(&destination, "destination", "d", "", "Destination directory in the vault (optional)")
}

// GetCommand returns the cp command for root command integration
func GetCommand() *cobra.Command {
	return cpCmd
}
