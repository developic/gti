package config

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	AppName = "gti"
)

var (
	ConfigDir  = filepath.Join(xdg.ConfigHome, AppName)
	DataDir    = filepath.Join(xdg.DataHome, AppName)
	CacheDir   = filepath.Join(xdg.CacheHome, AppName)
	ConfigFile = filepath.Join(ConfigDir, "config.toml")
)
