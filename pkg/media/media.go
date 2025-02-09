package media

import (
	"fmt"
	"github.com/sebastianappelberg/disk/pkg/config"
	"github.com/sebastianappelberg/disk/pkg/storage"
	"github.com/sebastianappelberg/disk/pkg/torrents"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type Type int

const (
	Unknown Type = iota
	Series
	Movie
)

type Media struct {
	Title             string
	Base              string
	Path              string
	Season            int
	Year              int
	Size              int64
	ModTime           time.Time
	Type              Type
	AvailabilityScore float64
}

func (m Media) GetPath() string {
	// We assume that the entire folder that the file resides in should be removed.
	return m.Base
}

func (m Media) GetPaths() []string {
	return []string{m.GetPath()}
}

func (m Media) String() string {
	if m.Type == Movie {
		return fmt.Sprintf("%s %d", m.Title, m.Year)
	}
	if m.Type == Series {
		return fmt.Sprintf("%s season %d", m.Title, m.Season)
	}
	return fmt.Sprintf("%s", m.Title)
}

type Analyzer struct {
	walker                     *storage.FileWalker[Media]
	sizeCalculator             *storage.SizeCalculator
	availabilityScoreThreshold float64
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		walker: storage.NewFileWalker[Media](
			storage.WithDecisionFilter[Media](decisionFilter),
			storage.WithMapper(contentMapper),
		),
		sizeCalculator:             storage.NewSizeCalculator(),
		availabilityScoreThreshold: 15,
	}
}

func decisionFilter(file storage.File) storage.FilterDecision {
	if config.ClutterFolders[file.Name] || config.UnsafeFolders[file.Name] {
		return storage.Skip
	}
	if isMediaFile(file.Name) {
		return storage.Include
	}
	return storage.Continue
}

func contentMapper(file storage.File, siblings []os.DirEntry) Media {
	content := Media{
		Size:    file.Size,
		ModTime: file.ModTime,
		Base:    file.Base,
		Path:    file.GetPath(),
	}
	inSeasonFolder := hasMultipleMediaFiles(siblings)

	parsedPath := ParsePath(content.Path)
	if parsedPath.Season != 0 {
		content.Title = parsedPath.Title
		content.Season = parsedPath.Season
		content.Type = Series
		return content
	}

	torrentInfo, err := torrents.ParseName(file.Name)
	if err == nil && content.Season == 0 {
		// If we were unable to find season by parsing the path, try with the torrentInfo.
		content.Title = torrentInfo.Title
		content.Season = torrentInfo.Season
		content.Year = torrentInfo.Year
	}
	if content.Season == 0 {
		if inSeasonFolder {
			// We're inside a season folder, but we were unable to find the season so we assume it is season 1.
			content.Season = 1
			content.Type = Series
		} else {
			if torrentInfo != nil && torrentInfo.Year == 0 {
				parse, err := torrents.ParseName(file.Base)
				if err == nil {
					content.Title = parse.Title
					content.Year = parse.Year
					content.Path = file.Base
				}
			}
			// If we're unable to find season, and it's only one file we assume it's a movie.
			content.Type = Movie
		}
	}
	return content
}

// hasMultipleMediaFiles checks if there is more than one media file
// among the provided directory entries.
func hasMultipleMediaFiles(entries []os.DirEntry) bool {
	mediaCount := 0
	for _, entry := range entries {
		if isMediaFile(entry.Name()) {
			mediaCount++
			if mediaCount > 1 {
				return true
			}
		}
	}
	return false
}

// Analyze returns a sorted list of candidates to delete.
func (a *Analyzer) Analyze(root string) []Media {
	contentCh := a.walker.GetFiles(root)

	seen := make(map[string]bool)
	var contents []Media
	for content := range contentCh {
		key := content.Title + strconv.Itoa(content.Season) + strconv.Itoa(content.Year)
		if !seen[key] {
			contents = append(contents, content)
			seen[key] = true
		}
	}

	contentWithAvailability := CheckAvailability(contents)
	var result []Media
	for _, content := range contentWithAvailability {
		if content.AvailabilityScore > a.availabilityScoreThreshold {
			if content.Type == Series {
				content.Size = a.sizeCalculator.GetSize(content.GetPath())
			}
			result = append(result, content)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].String() < result[j].String()
	})
	return result
}

// isMediaFile checks if the file is one of mp4, mkv, avi, mov, flv, wmv, webm, mp3, wav or flac.
// It doesn't deal with mixed-case filepath extensions for example: wAv. The reason being that
// this function is called a lot of times in a performance sensitive section and lower-casing
// the input reduces the performance of this function by 4x.
func isMediaFile(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".webm", ".flv", ".mp3", ".wav", ".flac",
		".MP4", ".MKV", ".AVI", ".MOV", ".WMV", ".WEBM", ".FLV", ".MP3", ".WAV", ".FLAC":
		return true
	}
	return false
}
