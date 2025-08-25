package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/imnotjames/micro-nas-lcd/internal/lcd"
	"github.com/imnotjames/micro-nas-lcd/internal/stats"
	"github.com/spf13/cobra"
)

func fmtKeyVal(key string, text string, maxLineLength uint8) string {
	if len(text) > int(maxLineLength)-4 {
		return fmt.Sprintf("%3s %12s", strings.ToUpper(key), text[:maxLineLength-4])
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
* Disk Information
* Disk Utilization
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

		interval, err := cmd.Flags().GetDuration("interval")
		if err != nil {
			panic(err)
		}

		diskDeviceNames, err := cmd.Flags().GetStringArray("disks")

		dev, err := lcd.NewDevice(address, columns, rows)
		if err != nil {
			panic(err)
		}
		defer dev.MustClose()

		if len(diskDeviceNames) == 0 {
			diskDeviceNames, err = stats.GetDisks()
			if err != nil {
				panic(err)
			}
		}

		for {
			hostText, _ := stats.GetHost()
			uptimeText, _ := stats.GetUptime()
			mustUpdateText(
				dev,
				fmtKeyVal("HST", hostText, columns),
				fmtKeyVal("UPT", uptimeText, columns),
			)

			time.Sleep(interval)

			memText, _ := stats.GetMemoryUtilization()
			swapText, _ := stats.GetSwapUtilization()
			mustUpdateText(
				dev,
				fmtKeyVal("MEM", memText, columns),
				fmtKeyVal("SWP", swapText, columns),
			)
			time.Sleep(interval)

			cpuText, _ := stats.GetCpuUtilization()
			loadText, _ := stats.GetLoad()
			mustUpdateText(
				dev,
				fmtKeyVal("CPU", cpuText, columns),
				loadText,
			)
			time.Sleep(interval)

			totalTransmit, _ := stats.GetTotalTransmit()
			totalReceive, _ := stats.GetTotalReceive()
			mustUpdateText(
				dev,
				fmtKeyVal("TTX", totalTransmit, columns),
				fmtKeyVal("TRX", totalReceive, columns),
			)
			time.Sleep(interval)

			connectionText, _ := stats.GetConnectionStatus("wlan0", "eth0")
			localIPText, _ := stats.GetLocalIP("wlan0", "eth0")
			mustUpdateText(
				dev,
				connectionText,
				localIPText,
			)
			time.Sleep(interval)

			for _, deviceName := range diskDeviceNames {
				diskInfoText, _ := stats.GetDiskInfo(deviceName)
				diskUtilizationText, _ := stats.GetDiskUtilization(deviceName)

				mustUpdateText(
					dev,
					diskInfoText,
					diskUtilizationText,
				)
				time.Sleep(interval)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)

	displayCmd.Flags().DurationP("interval", "i", 3*time.Second, "Interval between pages")
	displayCmd.Flags().StringArrayP("disks", "d", []string{}, "Disks to show")
}
