package cmd

import (
	"bufio"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gti/src/assets"
	"gti/src/internal/config"
)

var (
	listFlag    bool
	setFlag     string
	previewFlag string
)

var themeCmd = &cobra.Command{
	Use:   "theme [flags]",
	Short: "manage color themes for the typing interface",
	Long: `usage: gti theme [flags]

flags:
  --list              list all available themes (built-in and custom)
  --set <name>        set the active theme
  --preview <name>    preview a theme's colors without activating it`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		if listFlag {
			fmt.Println("Available themes:")
			for _, themeName := range getAvailableThemeNames() {
				fmt.Printf("  [✓] %s\n", themeName)
			}

		} else if setFlag != "" {
			if !isThemeAvailable(cfg, setFlag) {
				fmt.Printf("[ERROR] Theme '%s' is not available. Use --list to see available themes.\n", setFlag)
				return
			}

			cfg.Theme.Active = setFlag
			cfg.Theme.Colors = getThemeColors(setFlag)
			if err := config.SaveConfig(); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
			} else {
				fmt.Printf("[SUCCESS] Theme set to: %s\n", setFlag)
			}
		} else if previewFlag != "" {
			if !isThemeAvailable(cfg, previewFlag) {
				fmt.Printf("[ERROR] Theme '%s' is not available. Use --list to see available themes.\n", previewFlag)
				return
			}

			previewTheme(previewFlag)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	themeCmd.Flags().BoolVar(&listFlag, "list", false, "list all available themes (built-in and custom)")
	themeCmd.Flags().StringVar(&setFlag, "set", "", "set the active theme")
	themeCmd.Flags().StringVar(&previewFlag, "preview", "", "preview a theme's colors without activating it")
}

func isThemeAvailable(cfg *config.Config, themeName string) bool {
	themes := loadAvailableThemes()
	_, exists := themes[themeName]
	return exists
}

func previewTheme(themeName string) {
	themeColors := getThemeColors(themeName)

	fmt.Printf("Theme: %s\n", themeName)
	fmt.Printf("Background:     %s\n", themeColors.Background)
	fmt.Printf("Text Primary:   %s\n", themeColors.TextPrimary)
	fmt.Printf("Text Secondary: %s\n", themeColors.TextSecondary)
	fmt.Printf("Correct:        %s\n", themeColors.Correct)
	fmt.Printf("Incorrect:      %s\n", themeColors.Incorrect)
	fmt.Printf("Current:        %s\n", themeColors.Current)
	fmt.Printf("Pending:        %s\n", themeColors.Pending)
	fmt.Printf("Word Highlight: %s\n", themeColors.WordHighlight)
	fmt.Printf("Accent:         %s\n", themeColors.Accent)
	fmt.Printf("Border:         %s\n", themeColors.Border)
	fmt.Printf("Status Bar:     %s\n", themeColors.StatusBar)
}

func loadAvailableThemes() map[string]config.ThemeColorsConfig {
	themes := make(map[string]config.ThemeColorsConfig)

	entries, err := assets.Themes.ReadDir("themes")
	if err != nil {
		return themes
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		themeName := entry.Name()
		themePath := "themes/" + themeName

		if colors, err := loadThemeFromEmbeddedFile(themePath); err == nil {
			themes[themeName] = colors
		}
	}

	return themes
}

func loadThemeFromEmbeddedFile(filePath string) (config.ThemeColorsConfig, error) {
	data, err := assets.Themes.ReadFile(filePath)
	if err != nil {
		return config.ThemeColorsConfig{}, err
	}

	var colors config.ThemeColorsConfig
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.Contains(line, ":") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(strings.ToLower(strings.ReplaceAll(parts[0], " ", "_")))
		value := strings.TrimSpace(parts[1])

		switch key {
		case "background":
			colors.Background = value
		case "text_primary":
			colors.TextPrimary = value
		case "text_secondary":
			colors.TextSecondary = value
		case "correct":
			colors.Correct = value
		case "incorrect":
			colors.Incorrect = value
		case "current":
			colors.Current = value
		case "pending":
			colors.Pending = value
		case "word_highlight":
			colors.WordHighlight = value
		case "accent":
			colors.Accent = value
		case "border":
			colors.Border = value
		case "status_bar":
			colors.StatusBar = value
		}
	}

	return colors, scanner.Err()
}



func getThemeColors(themeName string) config.ThemeColorsConfig {
	themes := loadAvailableThemes()
	if colors, exists := themes[themeName]; exists {
		return colors
	}
	if defaultColors, exists := themes["default"]; exists {
		return defaultColors
	}
	return config.ThemeColorsConfig{
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
	}
}

func getAvailableThemeNames() []string {
	themes := loadAvailableThemes()
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
