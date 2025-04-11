// !! ---------------
// !! PAS TERMINER !!
// !! ---------------

package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/coyls/obs-cli/internal/config"
	"github.com/coyls/obs-cli/internal/logger"
)

var ExtractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract the latest backup from USB key",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runExtract()
	},
}

func init() {
	ArchiveCmd.AddCommand(ExtractCmd)
}

func runExtract() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Validate archive specific configuration
	if cfg.Config.Archive.UsbPath == "" {
		return fmt.Errorf("USB path for archive is required. Please set 'config.archive.usb_path' in your configuration")
	}
	if cfg.Config.Archive.ExtractPath == "" {
		return fmt.Errorf("extract path is required. Please set 'config.archive.extract_path' in your configuration")
	}

	logger.PrintHeader("Extract Obsidian Vaults Backup")

	// Check USB path
	usbPath := cfg.Config.Archive.UsbPath
	if _, err := os.Stat(usbPath); os.IsNotExist(err) {
		logger.Error("USB key is not connected at path: %s", usbPath)
		return fmt.Errorf("usb key not found at configured path")
	}

	logger.Info("USB key is connected at: %s", usbPath)

	// Find latest backup
	backupFile, err := findPreviousBackup(usbPath)
	if err != nil {
		return fmt.Errorf("failed to find backup: %w", err)
	}

	if backupFile == "" {
		logger.Error("No backup found on USB key")
		return fmt.Errorf("no backup found")
	}

	logger.Info("Found backup: %s", filepath.Base(backupFile))

	// Check if backup directory exists
	backupDir := cfg.Config.Archive.ExtractPath
	if _, err := os.Stat(backupDir); err == nil {
		logger.Info("Backup directory already exists: %s", backupDir)
		fmt.Print("Do you want to delete the existing directory and extract the backup? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" {
			logger.Info("Operation cancelled")
			return nil
		}

		logger.Info("Deleting existing directory...")
		if err := os.RemoveAll(backupDir); err != nil {
			return fmt.Errorf("failed to delete existing directory: %w", err)
		}
	}

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	logger.Info("Extracting backup to: %s", backupDir)

	// Extract backup
	if err := extractBackup(backupFile, backupDir); err != nil {
		return fmt.Errorf("failed to extract backup: %w", err)
	}

	logger.Success("Backup extracted successfully!")
	return nil
}

func extractBackup(backupFile, targetDir string) error {
	file, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		targetPath := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
			}

			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
			outFile.Close()
		}
	}

	return nil
}
