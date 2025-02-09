package media

import (
	"fmt"
	"github.com/sebastianappelberg/disk/pkg/torrents"
	"math"
	"sort"
	"strconv"
	"sync"
)

// CheckAvailability analyzes the availability of a given list of content.
// Availability as defined by how easy it'd be to get a hold of the content.
func CheckAvailability(content []Media) []Media {
	var wg sync.WaitGroup
	wg.Add(len(content))
	sem := make(chan struct{}, 20)
	result := make(chan Media)

	for _, c := range content {
		go func(ct Media) {
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()
			defer func() { result <- ct }()
			query := ct.String()
			torrentsResult, err := torrents.Search(query)
			if err != nil {
				fmt.Printf("Error searching torrents: %v\n", err)
				return
			}
			total := 0
			maxTorrentsLength := 3.0
			if ct.Type == Series {
				maxTorrentsLength = 10
			}
			torrentsLength := int(math.Min(maxTorrentsLength, float64(len(torrentsResult))))
			for _, torrent := range torrentsResult[:torrentsLength] {
				seeders, err := strconv.Atoi(torrent.Seeders)
				if err == nil {
					total += seeders
				}
			}
			average := float64(total) / float64(torrentsLength)
			ct.AvailabilityScore = average
		}(c)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	i := 0
	for res := range result {
		// Re-use the input slice instead of allocating a new one.
		content[i] = res
		i++
	}
	sort.Slice(content, func(i, j int) bool {
		return content[i].AvailabilityScore > content[j].AvailabilityScore
	})
	return content
}
