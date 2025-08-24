package cmd

import (
	"github.com/imnotjames/micro-nas-lcd/internal/lcd"
	"github.com/spf13/cobra"
)

var setTextCmd = &cobra.Command{
	Use:   "set-text",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		dev, err := lcd.NewDevice(0x20, 16, 2)
		if err != nil {
			panic(err)
		}
		defer dev.MustClose()

		if err := dev.UpdateText(args...); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setTextCmd)
}
