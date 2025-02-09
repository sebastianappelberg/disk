//go:build windows

package storage

import (
	"fmt"
	"golang.org/x/sys/windows"
)

// GetAvailableDisks lists all mounted disks.
func GetAvailableDisks() ([]string, error) {
	var disks []string
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
	return disks, nil
}
