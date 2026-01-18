package filesystem

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

type DiskStatus struct {
	Device      string // e.g., /dev/nvme0n1p3
	Mountpoint  string // e.g., /home
	Fstype      string // e.g., ext4, btrfs
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func GetAllDisks() ([]DiskStatus, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	var statuses []DiskStatus

	validTypes := map[string]bool{
		"ext4":    true,
		"btrfs":   true,
		"vfat":    true,
		"xfs":     true,
		"ntfs":    true,
		"fuseblk": true,
		"apfs":    true,
		"hfs+":    true,
	}

	for _, p := range partitions {
		if !validTypes[strings.ToLower(p.Fstype)] {
			continue
		}

		if shouldSkip(p.Fstype, p.Mountpoint) {
			continue
		}

		usage, err := disk.Usage(p.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue
		}

		statuses = append(statuses, DiskStatus{
			Device:      p.Device,
			Mountpoint:  p.Mountpoint,
			Fstype:      p.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	return statuses, nil
}

func shouldSkip(fstype, mountpoint string) bool {
	skipTypes := map[string]bool{
		"tmpfs": true, "devtmpfs": true, "devfs": true, "iso9660": true,
		"overlay": true, "squashfs": true, "proc": true, "sysfs": true,
	}

	if skipTypes[fstype] {
		return true
	}

	prefixes := []string{"/var/lib/docker", "/var/lib/flatpak", "/snap"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(mountpoint, prefix) {
			return true
		}
	}

	return false
}
