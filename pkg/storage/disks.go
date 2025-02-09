package storage

import (
	"bufio"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"runtime"
	"strings"
)

// GetAvailableDisks lists all mounted disks (cross-platform)
func GetAvailableDisks() ([]string, error) {
	var disks []string

	switch runtime.GOOS {
	case "windows":
		// Get bitmask of available drives
		drivesBitmask, err := windows.GetLogicalDrives()
		if err != nil {
			return nil, err
		}
		// Convert bitmask to drive letters
		for i := 0; i < 26; i++ {
			if drivesBitmask&(1<<uint(i)) != 0 {
				disks = append(disks, fmt.Sprintf("%c:\\", 'A'+i))
			}
		}
	case "linux", "darwin":
		// Open /proc/mounts to read mounted filesystems
		file, err := os.Open("/proc/mounts")
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) > 1 {
				mountPoint := fields[1] // The second column in /proc/mounts is the mount point
				disks = append(disks, mountPoint)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return disks, nil
}
