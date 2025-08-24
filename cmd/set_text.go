package cmd

import (
	"github.com/imnotjames/micro-nas-lcd/internal/lcd"
	"github.com/spf13/cobra"
)

var setTextCmd = &cobra.Command{
	Use:   "set-text",
	Short: "Set the LCD Display Text",
	Long: `
This updates the text displayed on the micronas LCD
panel.  Each "argument" passed is a line for the text.

For example,

	micro-nas-lcd set-text "Hello World" "Be Happy"

will display two lines of text on the LCD panel.

If the text is longer than the LCD panel it will be truncated.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		address, err := cmd.Flags().GetUint16("address")
		if err != nil {
			panic(err)
		}

		columns, err := cmd.Flags().GetUint8("columns")
		if err != nil {
			panic(err)
		}

		rows, err := cmd.Flags().GetUint8("rows")
		if err != nil {
			panic(err)
		}

		dev, err := lcd.NewDevice(address, columns, rows)
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
