package media

import (
	"fmt"
	"testing"
)

func TestParsePath(t *testing.T) {
	paths := []string{
		`D:\Media\TV Shows\Game of Thrones\Game of Thrones Season 4\Game.of.Thrones.S04E01.HDTV.x264-KILLERS.mp4`,
		`F:\My Series\Breaking Bad\Season 02\Breaking.Bad.S02E05.720p.BluRay.x264.mkv`,
		`E:\Videos\Stranger Things\Stranger Things Season 3\Stranger.Things.S03E08.WEB-DL.x264.mkv`,
		`C:\Users\Public\TV\Friends\Friends - Season 5\Friends.S05E12.1080p.WEBRip.x264.mp4`,
		`X:\TV Collection\The Office (US)\The Office Season 6\The.Office.S06E04.480p.HDTV.x264.mp4`,
		`Z:\Media Drive\The Boys\The Boys S02E07.720p.WEBRip.x265.mkv`, // No season folder
	}

	for _, path := range paths {
		parsedPath := ParsePath(path)
		fmt.Println(parsedPath)
	}
}
