package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

func GenerateConfig() error {
	dirs := []string{ConfigDir, DataDir, CacheDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	cfg := DefaultConfig()
	file, err := os.Create(ConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(cfg)
}
