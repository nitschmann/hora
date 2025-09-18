package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	FileName = "config"
	// ConfigLookupPaths defines the paths where the configuration file is looked up
	LookupPaths = []string{
		"./.hora",
		"$HOME/.hora",
		"~/.hora",
	}

	cfg Config

	defaultDebug                = false
	defaultListLimit            = 50
	defaultListOrder            = "desc"
	defaultUseBackgroundTracker = true
)

type Config struct {
	// DatabaseDir specifies the directory where the SQLite database file is stored
	DatabaseDir string `mapstructure:"database_dir" yaml:"database_dir"`
	Debug       bool   `mapstructure:"debug" yaml:"debug"`
	ListLimit   int    `mapstructure:"list_limit" yaml:"list_limit"`
	ListOrder   string `mapstructure:"list_order" yaml:"list_order"`
	// UseBackgroundTracker enables or disables the background tracker feature, which checks screen locks and (un)pauses time tracking based on these (macOS only for now)
	UseBackgroundTracker bool `mapstructure:"use_background_tracker" yaml:"use_background_tracker"`
}

// Load loads the configuration from the specified file or default locations.
// It returns the loaded Config, the path to the used configuration file (if any), and an error if occurred.
func Load(configFile string) (*Config, string, error) {
	defaultDatabaseDir, err := getDefaultDatabaseDir()
	if err != nil {
		return nil, "", err
	}

	viper.SetDefault("database_dir", defaultDatabaseDir)
	viper.SetDefault("debug", defaultDebug)
	viper.SetDefault("list_limit", defaultListLimit)
	viper.SetDefault("list_order", defaultListOrder)
	viper.SetDefault("use_background_tracker", defaultUseBackgroundTracker)

	viper.SetConfigType("yaml")

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		for _, path := range LookupPaths {
			viper.AddConfigPath(os.ExpandEnv(path))
		}
		viper.SetConfigName(FileName)
	}

	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return nil, "", err
		}
	}

	err = viper.Unmarshal(&cfg)

	return &cfg, viper.ConfigFileUsed(), err
}

func getDefaultDatabaseDir() (string, error) {
	var databaseDir string

	switch runtime.GOOS {
	case "darwin":
		// Use ~/Library/Application Support on macOS
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		databaseDir = filepath.Join(homeDir, "Library", "Application Support", "hora")
	default:
		// Use ~/.local/share on Linux and other Unix-like systems
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		databaseDir = filepath.Join(homeDir, ".local", "share", "hora")
	}

	return databaseDir, nil
}
