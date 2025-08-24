package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/imnotjames/micro-nas-lcd/internal/lcd"
	"github.com/imnotjames/micro-nas-lcd/internal/stats"
	"github.com/spf13/cobra"
)

const MaxLineLength = 16
const ScreenDuration = 3 * time.Second

func fmtKeyVal(key string, text string) string {
	if len(text) > MaxLineLength-4 {
		return fmt.Sprintf("%3s %12s", strings.ToUpper(key), text[:MaxLineLength-4])
	} else {
		return fmt.Sprintf("%3s %12s", strings.ToUpper(key), text)
	}
}

func mustUpdateText(dev *lcd.AdafruitLCDDevice, text ...string) {
	err := dev.UpdateText(text...)
	if err != nil {
		panic(err)
	}
}

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Cycles the LCD through helpful information.",
	Long: `
This will cycle through a variety of helpful pieces of information
about the NAS.

This includes:

* Hostname
* Uptime
* Memory Utilization
* Swap Utilization
* CPU Utilization
* System Load
* Transmitted Bytes
* Received Bytes
* Network Status
* Current IP
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

		for {
			hostText, _ := stats.GetHost()
			uptimeText, _ := stats.GetUptime()
			mustUpdateText(
				dev,
				fmtKeyVal("HST", hostText),
				fmtKeyVal("UPT", uptimeText),
			)

			time.Sleep(ScreenDuration)

			memText, _ := stats.GetMemoryUtilization()
			swapText, _ := stats.GetSwapUtilization()
			mustUpdateText(
				dev,
				fmtKeyVal("MEM", memText),
				fmtKeyVal("SWP", swapText),
			)

			time.Sleep(ScreenDuration)

			cpuText, _ := stats.GetCpuUtilization()
			loadText, _ := stats.GetLoad()
			mustUpdateText(dev, fmtKeyVal("CPU", cpuText), loadText)

			time.Sleep(ScreenDuration)

			totalTransmit, _ := stats.GetTotalTransmit()
			totalReceive, _ := stats.GetTotalReceive()
			mustUpdateText(
				dev,
				fmtKeyVal("TTX", totalTransmit),
				fmtKeyVal("TRX", totalReceive),
			)

			time.Sleep(ScreenDuration)

			connectionText, _ := stats.GetConnectionStatus("wlan0", "eth0")
			localIPText, _ := stats.GetLocalIP("wlan0", "eth0")
			mustUpdateText(dev, connectionText, localIPText)

			time.Sleep(ScreenDuration)
		}
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)
}
