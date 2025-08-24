package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "micro-nas-lcd",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
