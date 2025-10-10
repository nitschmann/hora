package config

import (
	"os"
	"path"

	"github.com/spf13/viper"

	"github.com/nitschmann/hora/internal/apperror"
)

func CreateDefault(directory string, forceOverwrite bool) (string, error) {
	var err error

	if directory == "" {
		directory = LookupPaths[0]
		expanded, err := expandPath(directory)
		if err != nil {
			return "", apperror.NewConfigErrorWithoutMessage(err)
		}
		directory = expanded

		// create directory if it does not exist
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return "", apperror.NewConfigErrorWithoutMessage(err)
		}
	}

	// set defaults
	defaultDatabaseDir, err := getDefaultDatabaseDir()
	if err != nil {
		return "", apperror.NewDatabaseError(err)
	}
	viper.Set("database_dir", defaultDatabaseDir)
	viper.Set("debug", defaultDebug)
	viper.Set("list_limit", defaultListLimit)
	viper.Set("list_order", defaultListOrder)
	viper.Set("use_background_tracker", defaultUseBackgroundTracker)
	viper.Set("web_ui_port", defaultWebUIPort)
	viper.Set("background_tracker_auto_stop", defaultBackgroundTrackerAutoStop)
	viper.Set("background_tracker_auto_stop_after", defaultBackgroundTrackerAutoStopAfter)

	configFilepath := path.Join(directory, FileName)

	err = viper.SafeWriteConfigAs(configFilepath)
	if err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok && forceOverwrite {
			err = viper.WriteConfigAs(configFilepath)
			if err != nil {
				return "", apperror.NewConfigErrorWithoutMessage(err)
			}
		} else {
			return "", apperror.NewConfigErrorWithoutMessage(err)
		}
	}

	_, createdConfigFile, err := Load(configFilepath)
	if err != nil {
		return "", apperror.NewConfigErrorWithoutMessage(err)
	}

	return createdConfigFile, nil
}
