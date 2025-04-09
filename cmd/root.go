package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/MUlt1mate/geny/app"
)

var rootCmd = &cobra.Command{
	Use:   "geny",
	Short: "Geny can convert and run go:generate statements as multi-line commands in YAML format",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		geny = &app.Geny{}
	},
}

type flags struct {
	File                 string
	Input                string
	Output               string
	Package              string
	HideGeneratedComment bool
}

var (
	runFlags flags
	geny     *app.Geny
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
