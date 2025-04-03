package storage

import (
	"fmt"
	"testing"
)

func TestGetAvailableDisks(t *testing.T) {
	disks, err := GetAvailableDisks()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(disks)
}
