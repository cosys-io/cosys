package common

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"os"
)

func InitConfigs() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigName(".cli_configs")
	viper.AddConfigPath(dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func GetPathConfig(key string, checkExists bool) (string, error) {
	if !viper.InConfig(key) {
		return "", fmt.Errorf("configuration not found: %s", key)
	}
	config := viper.GetString(key)

	if checkExists {
		exists, err := pathExists(config)
		if err != nil {
			return "", err
		}
		if !exists {
			return "", fmt.Errorf("path does not exist: %s", config)
		}
	}

	return config, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
