package trash

import (
	"testing"
)

func TestMoveToTrash(t *testing.T) {
	err := Put("test.txt")
	if err != nil {
		t.Fatalf("failed to move file to trash: %v", err)
	}
}

func TestRestore(t *testing.T) {
	err := Restore("test.txt")
	if err != nil {
		t.Fatalf("failed to move file to trash: %v", err)
	}
}
