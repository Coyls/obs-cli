package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/coyls/obs-cli/internal/config"
	"github.com/coyls/obs-cli/internal/logger"
)

var ArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Manage Obsidian vaults backups",
}

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a backup of all Obsidian vaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCreate()
	},
}

func init() {
	ArchiveCmd.AddCommand(CreateCmd)
}

func runCreate() error {
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

	logger.PrintHeader("Backup Obsidian Vaults")

	// Check USB path
	usbPath := cfg.Config.Archive.UsbPath
	if _, err := os.Stat(usbPath); os.IsNotExist(err) {
		logger.Error("USB key is not connected at path: %s", usbPath)
		return fmt.Errorf("usb key not found at configured path")
	}

	logger.Info("USB key is connected at: %s", usbPath)

	// Remove previous backup
	previousBackup, err := findPreviousBackup(usbPath)
	if err != nil {
		return fmt.Errorf("failed to find previous backup: %w", err)
	}

	if previousBackup != "" {
		logger.Info("Removing previous backup: %s", filepath.Base(previousBackup))
		if err := os.Remove(previousBackup); err != nil {
			return fmt.Errorf("failed to remove previous backup: %w", err)
		}
		logger.Success("Previous backup removed")
	} else {
		logger.Info("No previous backup found")
	}

	// Check available space
	requiredSpace, err := calculateRequiredSpace(cfg.Config.Root)
	if err != nil {
		return fmt.Errorf("failed to calculate required space: %w", err)
	}

	availableSpace, err := getAvailableSpace(usbPath)
	if err != nil {
		return fmt.Errorf("failed to get available space: %w", err)
	}

	if availableSpace < requiredSpace {
		logger.Error("Insufficient space on USB key")
		logger.Info("Required space: %s", formatBytes(requiredSpace))
		logger.Info("Available space: %s", formatBytes(availableSpace))
		return fmt.Errorf("insufficient space")
	}

	logger.Info("Sufficient space on USB key")
	logger.Info("Required space: %s", formatBytes(requiredSpace))
	logger.Info("Available space: %s", formatBytes(availableSpace))

	// Create backup
	backupFile := filepath.Join(usbPath, fmt.Sprintf("backup-obsidian_%s.tar.gz", getTimestamp()))

	logger.Info("Creating backup of all vaults...")
	if err := createBackup(cfg.Config.Root, backupFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	logger.Info("Verifying archive integrity...")
	if err := verifyBackup(backupFile); err != nil {
		os.Remove(backupFile)
		return fmt.Errorf("backup verification failed: %w", err)
	}

	logger.Success("Backup of all vaults completed and verified!")
	return nil
}

func findPreviousBackup(usbPath string) (string, error) {
	pattern := filepath.Join(usbPath, "backup-obsidian_*.tar.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(matches) > 0 {
		return matches[0], nil
	}
	return "", nil
}

func calculateRequiredSpace(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	// Add 10% margin
	size = size + (size / 10)
	return size, err
}

func getAvailableSpace(path string) (int64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}
	return int64(stat.Bavail) * int64(stat.Bsize), nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getTimestamp() string {
	return time.Now().Format("2006-01-02_15-04-05")
}

// cleanFileName nettoie le nom de fichier pour le rendre compatible avec les systèmes de fichiers
func cleanFileName(name string) string {
	// Remplacer les caractères problématiques par des tirets
	replacer := strings.NewReplacer(
		"?", "-",
		"*", "-",
		"<", "-",
		">", "-",
		"\"", "-",
		"'", "-",
	)
	return replacer.Replace(name)
}

func createBackup(sourcePath, targetFile string) error {
	file, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Convertir le chemin source en chemin absolu
	sourcePath, err = filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Convertir le chemin en absolu
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
		}

		// Créer un chemin relatif sécurisé
		relPath, err := filepath.Rel(sourcePath, absPath)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}

		// Nettoyer le nom de fichier
		dir := filepath.Dir(relPath)
		base := filepath.Base(relPath)
		cleanBase := cleanFileName(base)
		cleanPath := filepath.Join(dir, cleanBase)

		// Normaliser le chemin pour l'archive
		cleanPath = filepath.ToSlash(cleanPath)

		header, err := tar.FileInfoHeader(info, cleanPath)
		if err != nil {
			return fmt.Errorf("failed to create tar header for %s: %w", path, err)
		}

		// Définir le nom du fichier dans l'archive
		header.Name = cleanPath

		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header for %s: %w", path, err)
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return fmt.Errorf("failed to write file %s to archive: %w", path, err)
			}
		}

		return nil
	})
}

func verifyBackup(backupFile string) error {
	file, err := os.Open(backupFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		_, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCommand retourne la commande parent archive avec ses sous-commandes
func GetCommand() *cobra.Command {
	return ArchiveCmd
}
