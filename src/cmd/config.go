package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gti/src/internal/config"
)

var (
	showFlag  bool
	resetFlag bool
)

var configCmd = &cobra.Command{
	Use:   "config [flags]",
	Short: "view and manage GTI configuration settings",
	Long: `usage: gti config [flags]

flags:
  --show        display current configuration values
  --reset       reset configuration to default settings`,
	Run: func(cmd *cobra.Command, args []string) {
		if showFlag {
			cfg := config.GetConfig()
			fmt.Printf("Config file: %s\n\n", config.ConfigFile)
			printTimedConfig(cfg.Timed)
			printThemeConfig(cfg.Theme)
			printHistoryConfig(cfg.History)
		} else if resetFlag {
			fmt.Println("Resetting config to defaults...")
			if err := config.GenerateConfig(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Config reset successfully.")
			}
		} else {
			cmd.Help()
		}
	},
}

func printTimedConfig(timed config.TimedConfig) {
	fmt.Println("Timed:")
	fmt.Printf("  Default Seconds: %d\n", timed.DefaultSeconds)
	fmt.Println()
}

func printThemeConfig(theme config.ThemeConfig) {
	fmt.Println("Theme:")
	fmt.Printf("  Active: %s\n", theme.Active)
	fmt.Println()
}

func printHistoryConfig(history config.HistoryConfig) {
	fmt.Println("History:")
	fmt.Printf("  Enabled: %t\n", history.Enabled)
	fmt.Printf("  File:    %s\n", history.File)
	fmt.Println()
}

func init() {
	configCmd.Flags().BoolVar(&showFlag, "show", false, "display current configuration values")
	configCmd.Flags().BoolVar(&resetFlag, "reset", false, "reset configuration to default settings")
}
