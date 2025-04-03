//go:build darwin

package storage

import (
	"fmt"
	"golang.org/x/sys/unix"
	"strings"
)

// GetAvailableDisks returns a list of mounted disks.
func GetAvailableDisks() ([]string, error) {
	var mounts []string

	// Get the list of mounted file systems
	const maxEntries = 256
	mntbuf := make([]unix.Statfs_t, maxEntries)

	// Get mounted file system stats
	n, err := unix.Getfsstat(mntbuf, unix.MNT_NOWAIT)
	if err != nil {
		return nil, fmt.Errorf("failed to get mounted file systems: %v", err)
	}

	// Iterate over each mount point
	for i := 0; i < n; i++ {
		mountPoint := unix.ByteSliceToString(mntbuf[i].Mntonname[:])

		if !strings.HasPrefix(mountPoint, "/dev") && !strings.HasPrefix(mountPoint, "/System") {
			mounts = append(mounts, mountPoint)
		}
	}

	return mounts, nil
}
