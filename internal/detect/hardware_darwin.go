//go:build darwin

package detect

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// Pre-compiled regexes for diskutil output parsing.
var (
	reDeviceID   = regexp.MustCompile(`Device Identifier:\s*(\S+)`)
	reDiskSize   = regexp.MustCompile(`Disk Size:\s*[\d.]+\s*\w+\s*\((\d+)\s*Bytes\)`)
	reRemovable  = regexp.MustCompile(`Removable Media:\s*(\w+)`)
	reProtocol   = regexp.MustCompile(`Protocol:\s*(.+)`)
	reVolumeUUID = regexp.MustCompile(`Volume UUID:\s*(\S+)`)
	reParentDisk = regexp.MustCompile(`^(disk\d+)s\d+$`)
)

// HardwareInfo contains device-level information about the card/reader.
// On macOS, CID (Card Identification register) is not accessible through
// USB card readers. See hardware_linux.go for CID support.
type HardwareInfo struct {
	// Device size from diskutil (may differ from filesystem size due to partitioning/formatting)
	DeviceBytes int64

	// Filesystem size (what Statfs reports)
	FilesystemBytes int64

	// Block size
	BlockSize int64

	// Device identifier (e.g., "disk4s1")
	DeviceID string

	// Volume UUID if available
	VolumeUUID string

	// Whether this is a removable device
	IsRemovable bool

	// Protocol (USB, SD Card, etc.)
	Protocol string
}

// GetHardwareInfo attempts to retrieve hardware information for the given mount path.
// On macOS with USB readers, this returns limited info (no CID access).
func GetHardwareInfo(mountPath string) (*HardwareInfo, error) {
	info := &HardwareInfo{}

	// Get filesystem stats
	var stat syscall.Statfs_t
	if err := syscall.Statfs(mountPath, &stat); err != nil {
		return nil, err
	}
	info.FilesystemBytes = int64(stat.Blocks) * int64(stat.Bsize)
	info.BlockSize = int64(stat.Bsize)

	// Try to get device info from diskutil
	deviceID, err := getDeviceID(mountPath)
	if err != nil {
		return info, nil // Return what we have
	}
	info.DeviceID = deviceID

	// Query diskutil for device properties
	diskInfo, err := getDiskUtilInfo(deviceID)
	if err == nil {
		info.DeviceBytes = diskInfo.TotalSize
		info.IsRemovable = diskInfo.Removable
		info.Protocol = diskInfo.Protocol
		info.VolumeUUID = diskInfo.VolumeUUID
	}

	return info, nil
}

type diskUtilInfo struct {
	TotalSize  int64
	Removable  bool
	Protocol   string
	VolumeUUID string
}

func getDeviceID(mountPath string) (string, error) {
	cmd := exec.Command("diskutil", "info", mountPath)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	matches := reDeviceID.FindStringSubmatch(string(out))
	if len(matches) >= 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("device identifier not found")
}

func getDiskUtilInfo(deviceID string) (*diskUtilInfo, error) {
	// Get parent disk (e.g., disk4s1 -> disk4)
	parentDisk := deviceID
	if m := reParentDisk.FindStringSubmatch(deviceID); len(m) >= 2 {
		parentDisk = m[1]
	}

	// Query the physical disk for size info
	cmd := exec.Command("diskutil", "info", parentDisk)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	info := &diskUtilInfo{}
	output := string(out)

	if m := reDiskSize.FindStringSubmatch(output); len(m) >= 2 {
		size, _ := strconv.ParseInt(m[1], 10, 64)
		info.TotalSize = size
	}

	if m := reRemovable.FindStringSubmatch(output); len(m) >= 2 {
		info.Removable = m[1] == "Yes" || m[1] == "Removable"
	}

	if m := reProtocol.FindStringSubmatch(output); len(m) >= 2 {
		info.Protocol = strings.TrimSpace(m[1])
	}

	// Get volume UUID from the partition device
	cmd = exec.Command("diskutil", "info", deviceID)
	out, _ = cmd.Output()
	if m := reVolumeUUID.FindStringSubmatch(string(out)); len(m) >= 2 {
		info.VolumeUUID = m[1]
	}

	return info, nil
}

// FormatHardwareInfo returns a formatted string with hardware details.
func FormatHardwareInfo(info *HardwareInfo) string {
	if info == nil {
		return "Hardware info unavailable"
	}

	var parts []string

	parts = append(parts, fmt.Sprintf("Device: %s", info.DeviceID))

	if info.Protocol != "" {
		parts = append(parts, fmt.Sprintf("Protocol: %s", info.Protocol))
	}

	if info.DeviceBytes > 0 {
		parts = append(parts, fmt.Sprintf("Raw Size: %s", FormatBytes(info.DeviceBytes)))
	}

	if info.FilesystemBytes > 0 {
		parts = append(parts, fmt.Sprintf("Filesystem: %s", FormatBytes(info.FilesystemBytes)))
	}

	if info.IsRemovable {
		parts = append(parts, "Removable: Yes")
	}

	if info.VolumeUUID != "" {
		parts = append(parts, fmt.Sprintf("UUID: %s", info.VolumeUUID))
	}

	parts = append(parts, "CID: Not available (USB reader)")

	return strings.Join(parts, "\n  ")
}
