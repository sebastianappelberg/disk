package torrents

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	torrents, err := Search("Some movie title")
	if err != nil {
		t.Error(err)
	}
	for _, torrent := range torrents {
		fmt.Println(torrent.Name)
	}
}
