package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	VaultPath     string `yaml:"vault_path"`
	DefaultMvPath string `yaml:"default_mv_path"`
	DefaultCpPath string `yaml:"default_cp_path"`
	ObsidianPath  string `yaml:"obsidian_path"`  // Path to .obsidian directory
	DefaultEditor string `yaml:"default_editor"` // Default editor to use (code, nano, vim)
}

// findObsidianVault finds the Obsidian vault path by looking for .obsidian directory
func findObsidianVault() (string, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Run find command to locate .obsidian directory
	cmd := exec.Command("find", homeDir, "-type", "d", "-name", ".obsidian")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to find Obsidian vault: %w", err)
	}

	// Split output into lines and get the first result
	paths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(paths) == 0 || paths[0] == "" {
		return "", fmt.Errorf("no Obsidian vault found")
	}

	// Get the parent directory of .obsidian
	vaultPath := filepath.Dir(paths[0])
	return vaultPath, nil
}

// findObsidianConfig finds the .obsidian directory path
func findObsidianConfig(vaultPath string) string {
	return filepath.Join(vaultPath, ".obsidian")
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	// Try to find Obsidian vault path
	vaultPath, err := findObsidianVault()
	if err != nil {
		// Fallback to environment variable if vault not found
		vaultPath = os.Getenv("OBSIDIAN_VAULTS_PATH")
	}

	// Get .obsidian path
	obsidianPath := findObsidianConfig(vaultPath)

	return &Config{
		VaultPath:     vaultPath,
		DefaultMvPath: "Archives", // Default directory for mv command
		DefaultCpPath: "Assets",   // Default directory for cp command
		ObsidianPath:  obsidianPath,
		DefaultEditor: "code", // Default editor is VS Code
	}
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Define config file path
	configPath := filepath.Join(homeDir, ".config", "obs-cli", "config.yaml")

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config if file doesn't exist
			config := DefaultConfig()
			if err := config.Save(configPath); err != nil {
				return nil, err
			}
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config file
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves the configuration to a file
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.VaultPath == "" {
		return fmt.Errorf("no Obsidian vault path found. please set OBSIDIAN_VAULTS_PATH environment variable or run 'obs-cli init'")
	}

	// Check if vault directory exists
	if _, err := os.Stat(c.VaultPath); os.IsNotExist(err) {
		return fmt.Errorf("obsidian vault directory not found: %s", c.VaultPath)
	}

	// Check if .obsidian directory exists
	if _, err := os.Stat(c.ObsidianPath); os.IsNotExist(err) {
		return fmt.Errorf("obsidian config directory not found: %s", c.ObsidianPath)
	}

	// Check if default mv path exists
	if c.DefaultMvPath != "" {
		defaultMvPath := filepath.Join(c.VaultPath, c.DefaultMvPath)
		if _, err := os.Stat(defaultMvPath); os.IsNotExist(err) {
			// Create default mv directory if it doesn't exist
			if err := os.MkdirAll(defaultMvPath, 0755); err != nil {
				return fmt.Errorf("failed to create default mv directory: %w", err)
			}
		}
	}

	// Check if default cp path exists
	if c.DefaultCpPath != "" {
		defaultCpPath := filepath.Join(c.VaultPath, c.DefaultCpPath)
		if _, err := os.Stat(defaultCpPath); os.IsNotExist(err) {
			// Create default cp directory if it doesn't exist
			if err := os.MkdirAll(defaultCpPath, 0755); err != nil {
				return fmt.Errorf("failed to create default cp directory: %w", err)
			}
		}
	}

	// Validate editor
	if c.DefaultEditor == "" {
		c.DefaultEditor = "code" // Set default editor if not specified
	}

	return nil
}
