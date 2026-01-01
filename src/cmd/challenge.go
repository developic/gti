package cmd

import (
	"gti/src/internal/app"

	"github.com/spf13/cobra"
)

var challengeCmd = &cobra.Command{
	Use:   "challenge",
	Short: "Start progressive challenge mode with levels",
	Long: `Start challenge mode with progressive difficulty levels and boss rounds.
Complete increasingly difficult typing challenges to unlock achievements.

EXAMPLES:
  gti challenge    # Start from current level

CONTROLS: Same as other modes
  During challenge:
    Type normally to meet requirements
    Tab/Enter    Submit when requirements met
    Ctrl+C       Quit challenge

  Results screen:
    Enter         Continue to next level
    Ctrl+R        Retry current level
    Ctrl+C        Quit to main menu

PROGRESS:
  Challenge progress is saved automatically
  Failed attempts don't reset progress`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.StartChallengeGame()
	},
}
