package config

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

const DefaultPracticeText = "Typing is not about speed alone, it is about accuracy, rhythm, and calm focus."

type Config struct {
	Display  DisplayConfig  `toml:"display"`
	Theme    ThemeConfig    `toml:"theme"`
	Timed    TimedConfig    `toml:"timed"`
	Language LanguageConfig `toml:"language"`
	Network  NetworkConfig  `toml:"network"`
	History  HistoryConfig  `toml:"history"`
}

type DisplayConfig struct {
	MaxWidth        int  `toml:"max_width"`
	CenterText      bool `toml:"center_text"`
	ShowProgressBar bool `toml:"show_progress_bar"`
	FPS             int  `toml:"fps"`
}

type ThemeConfig struct {
	Active string `toml:"active"`
	Colors ThemeColorsConfig
	Styles ThemeStylesConfig
}

type ThemeColorsConfig struct {
	Correct       string
	Incorrect     string
	Current       string
	Pending       string
	WordHighlight string
	Accent        string
	Border        string
	TextPrimary   string
	TextSecondary string
	Background    string
	StatusBar     string
}

type ThemeStylesConfig struct {
	UnderlineCurrent bool `toml:"underline_current"`
	DimPending       bool `toml:"dim_pending"`
	BoldResults      bool `toml:"bold_results"`
}

type TimedConfig struct {
	DefaultSeconds int `toml:"default_seconds"`
}

type LanguageConfig struct {
	Default string `toml:"default"`
}

type NetworkConfig struct {
	TimeoutMs int `toml:"timeout_ms"`
}

type HistoryConfig struct {
	Enabled bool   `toml:"enabled"`
	File    string `toml:"file"`
}

func DefaultConfig() *Config {
	return &Config{
		Display: DisplayConfig{
			MaxWidth:        80,
			CenterText:      true,
			ShowProgressBar: true,
			FPS:             60,
		},

		Theme: ThemeConfig{
			Active: "gruvbox",
			Colors: ThemeColorsConfig{
				Background:    "#000000",
				TextPrimary:   "#FFFFFF",
				TextSecondary: "#AAAAAA",
				Correct:       "#00C853",
				Incorrect:     "#FF1744",
				Current:       "#FFD600",
				Pending:       "#7C7C7C",
				Accent:        "#00AAFF",
				Border:        "#444444",
				StatusBar:     "#333333",
			},
			Styles: ThemeStylesConfig{
				UnderlineCurrent: true,
				DimPending:       true,
				BoldResults:      true,
			},
		},
		Timed: TimedConfig{
			DefaultSeconds: 30,
		},
		Language: LanguageConfig{
			Default: "english",
		},

		Network: NetworkConfig{
			TimeoutMs: 5000,
		},
		History: HistoryConfig{
			Enabled: true,
			File:    filepath.Join(xdg.DataHome, "gti", "history.jsonl"),
		},
	}
}
