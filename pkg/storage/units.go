package storage

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	MegaByte = KiloByte * KiloByte
	GigaByte = MegaByte * KiloByte

	byteSuffix     = "B"
	kiloByteSuffix = "kB"
	megaByteSuffix = "MB"
	gigaByteSuffix = "GB"
)

func FormatSize[integer Integer](size integer) string {
	if size == 0 {
		return "0"
	}
	if size < KiloByte {
		return fmt.Sprintf("%d%s", size, byteSuffix)
	}
	if size < MegaByte {
		return fmt.Sprintf("%d%s", size/KiloByte, kiloByteSuffix)
	}
	if size < GigaByte {
		return fmt.Sprintf("%d%s", size/MegaByte, megaByteSuffix)
	}
	return fmt.Sprintf("%d%s", size/GigaByte, gigaByteSuffix)
}

func ParseSize(size string) (int64, error) {
	if strings.HasSuffix(size, kiloByteSuffix) {
		n, err := splitToInt(size, kiloByteSuffix)
		if err != nil {
			return 0, err
		}
		return KiloByte * n, nil
	}
	if strings.HasSuffix(size, megaByteSuffix) {
		n, err := splitToInt(size, megaByteSuffix)
		if err != nil {
			return 0, err
		}
		return MegaByte * n, nil
	}
	if strings.HasSuffix(size, gigaByteSuffix) {
		n, err := splitToInt(size, gigaByteSuffix)
		if err != nil {
			return 0, err
		}
		return GigaByte * n, nil
	}
	// Important to check for 'B' last since 'B' would be caught by all other strings 'kB'.
	if strings.HasSuffix(size, byteSuffix) {
		n, err := splitToInt(size, byteSuffix)
		if err != nil {
			return 0, err
		}
		return n, nil
	}
	return 0, fmt.Errorf("invalid size format: %s", size)
}

func splitToInt(size string, unit string) (int64, error) {
	split := strings.Split(size, unit)
	atoi, err := strconv.Atoi(split[0])
	if atoi < 0 || err != nil {
		return 0, fmt.Errorf("invalid size format: %s", size)
	}
	return int64(atoi), nil
}

type Integer interface {
	int64 | uint64
}
