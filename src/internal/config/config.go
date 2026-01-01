package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

var globalConfig *Config

func InitConfig(configFile string) {
	if configFile != "" {
		ConfigFile = ExpandPath(configFile)
		ConfigDir = filepath.Dir(ConfigFile)
		CacheDir = filepath.Join(ConfigDir, "cache")
	}

	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		if err := GenerateConfig(); err != nil {
			panic("Failed to generate config: " + err.Error())
		}
	}

	if err := LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config file %s: %v\nUsing default configuration.\n", ConfigFile, err)
		globalConfig = DefaultConfig()
	}
}

func LoadConfig() error {
	file, err := os.Open(ConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg := DefaultConfig()
	_, err = toml.DecodeReader(file, cfg)
	if err != nil {
		return err
	}

	globalConfig = cfg
	return nil
}

func GetConfig() *Config {
	if globalConfig == nil {
		InitConfig("")
	}
	return globalConfig
}

func SaveConfig() error {
	file, err := os.Create(ConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(globalConfig)
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return strings.Replace(path, "~", home, 1)
	}
	return path
}
