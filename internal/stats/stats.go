package stats

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/sensors"
)

const ListDisksTimeout = 5 * time.Second
const DiskUsageTimeout = 5 * time.Second
const CpuUtilizationTimeout = 5 * time.Second
const NetInterfacesTimeout = 5 * time.Second

func fmtBytes(b uint64) string {
	return fmtBytesPrecision(b, 0, 1.0)
}

func fmtBytesPrecision(b uint64, precision uint8, threshold float64) string {
	format := fmt.Sprintf("%%.%df%%c", precision)

	const unit = 1000
	unitThreshold := uint64(unit * threshold)
	if b < unit {
		return fmt.Sprintf(format, float64(b), 'B')
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unitThreshold && exp < 6; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf(format, float64(b)/float64(div), "KMGTPE"[exp])
}

func fmtMemoryUtilization(used uint64, total uint64, usedPercent float64) string {
	return fmt.Sprintf(
		"%3s/%3s %3.0f%%",
		fmtBytes(used),
		fmtBytes(total),
		usedPercent,
	)
}

func GetMemoryUtilization() (string, error) {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	return fmtMemoryUtilization(virtualMemory.Used, virtualMemory.Total, virtualMemory.UsedPercent), nil
}

func GetSwapUtilization() (string, error) {
	swapMemory, err := mem.SwapMemory()
	if err != nil {
		return "", err
	}
	return fmtMemoryUtilization(swapMemory.Used, swapMemory.Total, swapMemory.UsedPercent), nil
}

func GetHost() (string, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return "", err
	}
	return hostInfo.Hostname, nil
}

func GetUptime() (string, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return "", err
	}

	return time.Duration(hostInfo.Uptime * uint64(time.Second)).String(), nil
}

func GetCpuUtilization() (string, error) {
	percents, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", err
	}

	maxPercent := 0.0
	for _, percent := range percents {
		maxPercent = math.Max(maxPercent, percent)
	}

	ctx, cancel := context.WithTimeout(context.Background(), CpuUtilizationTimeout)
	defer cancel()

	temperatures, err := sensors.TemperaturesWithContext(ctx)
	if err != nil {
		return "", err
	}
	maxTemperature := 0.0
	for _, temperature := range temperatures {
		if strings.HasPrefix(temperature.SensorKey, "coretemp_core_") {
			maxTemperature = math.Max(maxTemperature, temperature.Temperature)
		}
	}

	return fmt.Sprintf("%6.2f%% %3.0fC", maxPercent, maxTemperature), nil
}

func GetLoad() (string, error) {
	loadAverages, err := load.Avg()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.1f/%.1f/%.1f", loadAverages.Load1, loadAverages.Load5, loadAverages.Load15), nil
}

func GetTotalTransmit() (string, error) {
	counters, err := net.IOCounters(false)
	if err != nil {
		return "", err
	}
	total := uint64(0)
	for _, counter := range counters {
		total += counter.BytesSent
	}

	return fmt.Sprintf("%4s", fmtBytesPrecision(total, 2, 9)), nil
}

func GetTotalReceive() (string, error) {
	counters, err := net.IOCounters(false)
	if err != nil {
		return "", err
	}
	total := uint64(0)
	for _, counter := range counters {
		total += counter.BytesRecv
	}

	return fmt.Sprintf("%4s", fmtBytesPrecision(total, 2, 9)), nil
}

func getInterfaces(interfaceNames ...string) ([]net.InterfaceStat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), NetInterfacesTimeout)
	defer cancel()

	interfaces, err := net.InterfacesWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return slices.Collect(func(yield func(net.InterfaceStat) bool) {
		for _, iface := range interfaces {
			if slices.Contains(interfaceNames, iface.Name) {
				if !yield(iface) {
					return
				}
			}
		}
	}), nil
}

func GetConnectionStatus(interfaceNames ...string) (string, error) {
	interfaces, err := getInterfaces(interfaceNames...)
	if err != nil {
		return "", err
	}

	if len(interfaces) == 0 {
		return "", fmt.Errorf("no interface found")
	}

	iface := interfaces[0]
	if slices.Contains(iface.Flags, "up") {
		return fmt.Sprintf("%s CONNECTED", iface.Name), nil
	} else {
		return fmt.Sprintf("%s DISCONNECTED", iface.Name), nil
	}
}

func GetLocalIP(interfaceNames ...string) (string, error) {
	interfaces, err := getInterfaces(interfaceNames...)
	if err != nil {
		return "", err
	}

	if len(interfaces) == 0 {
		return "", fmt.Errorf("no interface found")
	}

	for _, iface := range interfaces {
		if len(iface.Addrs) > 0 {
			addr, _, _ := strings.Cut(iface.Addrs[0].Addr, "/")
			return addr, nil
		}
	}

	return "", nil
}

func getDeviceUsage(device string) (*disk.UsageStat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DiskUsageTimeout)
	defer cancel()

	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return nil, err
	}

	for _, partition := range partitions {
		if partition.Device != device {
			continue
		}

		usage, err := disk.UsageWithContext(ctx, partition.Mountpoint)
		if err != nil {
			return nil, err
		}

		return usage, nil
	}

	return nil, fmt.Errorf("no matching device")
}

func GetDisks() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ListDisksTimeout)
	defer cancel()

	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return nil, err
	}

	disks := make([]string, len(partitions))

	for i, partition := range partitions {
		disks[i] = partition.Device
	}

	slices.Sort(disks)
	return slices.Compact(disks), nil
}

func GetDiskInfo(device string) (string, error) {
	usage, err := getDeviceUsage(device)

	shortDeviceName := strings.TrimPrefix(device, "/dev/")

	if err != nil {
		return fmt.Sprintf("no disk: %s", shortDeviceName), err
	}

	fsType := usage.Fstype

	if fsType == "ext2/ext3" || fsType == "ext2" {
		fsType = "ext"
	}

	return fmt.Sprintf("%s %s", shortDeviceName, fsType), nil
}

func GetDiskUtilization(device string) (string, error) {
	usage, err := getDeviceUsage(device)
	if err != nil {
		return "", err
	}

	total := fmtBytes(usage.Total)

	return fmt.Sprintf("%s %.2f%%",
		total,
		usage.UsedPercent,
	), nil
}
