package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print GTI version",
	Long:  "Print the current GTI version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gti %s\n", Version)
	},
}
