//go:build linux

package storage

import (
	"bufio"
	"os"
	"strings"
)

// GetAvailableDisks lists all mounted disks.
func GetAvailableDisks() ([]string, error) {
	var disks []string
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
	return disks, nil
}
