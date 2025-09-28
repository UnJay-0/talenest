package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	SQLitePath string `mapstructure:"sqlite_path"`
	DuckDBpath string `mapstructure:"sqlite_path"`
}

func GetAppDataPath(filename string) string {
	var baseDir string

	// TODO: find another way to identify the OS
	if appData := os.Getenv("APPDATA"); appData != "" {
		// Windows
		baseDir = filepath.Join(appData, "Talenest")
	} else if home := os.Getenv("HOME"); home != "" {
		// TODO: macos to change base dir
		// macOS -> $HOME/Library/Application Support/
		// Linux
		baseDir = filepath.Join(home, ".local", "share", "Talenest")
	} else {
		baseDir = "."
	}

	_ = os.MkdirAll(baseDir, 0700)
	return filepath.Join(baseDir, filename)
}

func LoadConfig() *Config {
	appDir := GetAppDataPath("")

	// create app data directories
	InitProgramDirectories(appDir)

	configPath := filepath.Join(appDir, "config")

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)

	viper.SetDefault("sqlite_path", filepath.Join(appDir, "data", "talenest.db"))
	viper.SetDefault("duckdb_path", filepath.Join(appDir, "data", "talenest_analytics.duckdb"))

	if err := viper.ReadInConfig(); err != nil {
		// If missing, write defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				log.Fatalf("Failed to write default config: %v", err)
			}
		} else {
			log.Fatalf("Failed to read config: %v", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	return &cfg
}

// TODO: Check if it's a proper way
func InitProgramDirectories(appDir string) {
	// TODO: standardize directories name
	directories := []string{"config", "data"}
	for _, directory := range directories {
		if err := os.MkdirAll(filepath.Join(appDir, directory), 0700); err != nil {
			log.Fatalf("Failed to create app dir: %v", err)
		}
	}
}
