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

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if _, err := os.Stat(source); os.IsNotExist(err) {
		logger.Error("Source file not found: %s", source)
		return fmt.Errorf("source file not found: %s", source)
	}

	fmt.Printf("Default Vault: %s\n", cfg.Config.DefaultVault)
	for key := range cfg.Config.Vaults {
		fmt.Printf("Vault trouvé: %s\n", key)
	}

	vaultConfig, exists := cfg.GetVaultConfig(cfg.Config.DefaultVault)
	if !exists {
		logger.Error("Default vault configuration not found")
		return fmt.Errorf("default vault configuration not found")
	}

	if destination == "" {
		if vaultConfig.Commands.Cp.DefaultTargetPath == "" {
			logger.Error("No destination specified and no default path configured")
			return fmt.Errorf("no destination specified and no default path configured")
		}
		destination = vaultConfig.Commands.Cp.DefaultTargetPath
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

	srcFile, err := os.Open(source)
	if err != nil {
		logger.Error("Failed to open source file: %s", err.Error())
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(finalDest)
	if err != nil {
		logger.Error("Failed to create destination file: %s", err.Error())
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

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

func GetCommand() *cobra.Command {
	return cpCmd
}
