package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type VaultConfig struct {
	VaultPath string `mapstructure:"vault_path"`
	Commands  struct {
		Cp struct {
			DefaultTargetPath string `mapstructure:"default_target_path"`
		} `mapstructure:"cp"`
		Mv struct {
			DefaultTargetPath string `mapstructure:"default_target_path"`
		} `mapstructure:"mv"`
	} `mapstructure:"commands"`
}

type Config struct {
	Config struct {
		DefaultEditor string                  `mapstructure:"default_editor"`
		DefaultVault  string                  `mapstructure:"default_vault"`
		Root          string                  `mapstructure:"root"`
		Vaults        map[string]*VaultConfig `mapstructure:"vaults"`
		Archive       struct {
			UsbPath     string `mapstructure:"usb_path"`
			ExtractPath string `mapstructure:"extract_path"`
		} `mapstructure:"archive"`
	} `mapstructure:"config"`
}

func LoadConfig() (*Config, error) {

	path := os.Getenv("OBS_CLI_CONFIG")
	if path == "" {
		return nil, fmt.Errorf("OBS_CLI_CONFIG is not set")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", path)
	}

	v := viper.New()

	v.SetConfigType("yaml")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	// debugConfig(c)

	if c.Config.Root == "" {
		return fmt.Errorf("vault root path is required")
	}
	if c.Config.DefaultVault == "" {
		return fmt.Errorf("default vault name is required")
	}
	if _, exists := c.GetVaultConfig(c.Config.DefaultVault); !exists {
		return fmt.Errorf("configuration for default vault '%s' not found", c.Config.DefaultVault)
	}
	return nil
}

// GetVaultConfig retourne la configuration d'un vault de manière insensible à la casse
func (c *Config) GetVaultConfig(vaultName string) (*VaultConfig, bool) {
	// Chercher d'abord avec la clé exacte
	if config, exists := c.Config.Vaults[vaultName]; exists {
		return config, true
	}

	// Si non trouvé, chercher sans tenir compte de la casse
	lowerName := strings.ToLower(vaultName)
	for key, config := range c.Config.Vaults {
		if strings.ToLower(key) == lowerName {
			return config, true
		}
	}
	return nil, false
}

// func debugConfig(cfg *Config) {
// 	vaultConfig, exists := cfg.GetVaultConfig("Coyls")
// 	fmt.Printf("Config for 'Coyls': %+v (exists: %v)\n", vaultConfig, exists)

// 	vaultConfig, exists = cfg.GetVaultConfig(cfg.Config.DefaultVault)
// 	fmt.Printf("Config for default vault '%s': %+v (exists: %v)\n", cfg.Config.DefaultVault, vaultConfig, exists)
// }
