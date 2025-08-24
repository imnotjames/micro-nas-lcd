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

func init() {
	rootCmd.PersistentFlags().Uint16P("address", "a", 0x20, "I2C Address of the Adafruit LCD Backpack")
	rootCmd.PersistentFlags().Uint8P("columns", "c", 16, "Columns available on the LCD Panel")
	rootCmd.PersistentFlags().Uint8P("rows", "r", 2, "Rows available on the LCD Panel")
}
