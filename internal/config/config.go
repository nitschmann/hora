package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/spf13/viper"
)

var (
	FileName = "config.yaml"
	// ConfigLookupPaths defines the paths where the configuration file is looked up
	LookupPaths = []string{
		"~/.hora",
		"./.hora",
	}

	cfg Config

	defaultDebug                = false
	defaultListLimit            = 50
	defaultListOrder            = "desc"
	defaultUseBackgroundTracker = true
	defaultWebUIPort            = 8080
)

type Config struct {
	// DatabaseDir specifies the directory where the SQLite database file is stored
	DatabaseDir string `mapstructure:"database_dir" yaml:"database_dir"`
	Debug       bool   `mapstructure:"debug" yaml:"debug"`
	ListLimit   int    `mapstructure:"list_limit" yaml:"list_limit" validate:"gte=1"`
	ListOrder   string `mapstructure:"list_order" yaml:"list_order" validate:"oneof=asc desc"`
	// UseBackgroundTracker enables or disables the background tracker feature, which checks screen locks and (un)pauses time tracking based on these (macOS only for now)
	UseBackgroundTracker bool `mapstructure:"use_background_tracker" yaml:"use_background_tracker"`
	WebUIPort            int  `mapstructure:"web_ui_port" yaml:"web_ui_port" validate:"gte=1,lte=65535"`
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
	viper.SetDefault("web_ui_port", defaultWebUIPort)

	viper.SetConfigType("yaml")

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		for _, path := range LookupPaths {
			expandedPath, err := expandPath(path)
			if err != nil {
				return nil, "", err
			}
			viper.AddConfigPath(expandedPath)
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
	if err != nil {
		return nil, "", err
	}

	err = validateConfig(&cfg)
	if err != nil {
		return nil, "", err
	}

	return &cfg, viper.ConfigFileUsed(), err
}

// expandPath expands the ~ in the given path to the user's home directory
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			return home, nil
		}

		return filepath.Join(home, path[1:]), nil
	}

	return path, nil
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

func validateConfig(cfg *Config) error {
	en := en.New()
	uni := ut.New(en, en)
	validationTranslator, _ := uni.GetTranslator("en")

	validate := validator.New()
	err := en_translations.RegisterDefaultTranslations(validate, validationTranslator)
	if err != nil {
		return err
	}

	err = validate.Struct(cfg)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorCount := len(validationErrors)
		validationErrorMessages := make([]string, len(validationErrors))

		for _, e := range validationErrors {
			validationErrorMessages = append(validationErrorMessages, e.Translate(validationTranslator))
		}

		if errorCount == 1 {
			return fmt.Errorf("configuration validation error\n%s\n", strings.Join(validationErrorMessages, "\n"))
		} else {
			return fmt.Errorf("configuration validation errors\n%s\n", strings.Join(validationErrorMessages, "\n"))
		}
	}

	return nil
}
