package cmd

import (
	"gti/src/internal/config"
	"gti/src/internal/session"
	"gti/src/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var quoteCount int

var quoteCmd = &cobra.Command{
	Use:   "quote [options]",
	Short: "start quote typing mode (default)",
	Long: `usage: gti quote [options]

options:
  -n, --count <num>    number of quotes to type (default: 2)
  -h, --help           display help information`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetConfig()
		quoteList := session.FetchMultipleQuotes(cfg, quoteCount)
		sess := session.NewSessionWithQuotes(cfg, quoteList)
		model := tui.NewModelWithSession(cfg, sess)
		p := tea.NewProgram(model, tea.WithAltScreen())
		_, err := p.Run()
		return err
	},
}

func init() {
	quoteCmd.Flags().IntVarP(&quoteCount, "count", "n", 2, "number of quotes to type")
}
