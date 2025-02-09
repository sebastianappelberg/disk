package clean

import (
	"github.com/sebastianappelberg/disk/pkg/clutter"
	"github.com/sebastianappelberg/disk/pkg/config"
	"github.com/sebastianappelberg/disk/pkg/games"
	"github.com/sebastianappelberg/disk/pkg/media"
	"github.com/sebastianappelberg/disk/pkg/trash"
	"time"
)

type Args struct {
	Root        string
	MinAge      int
	MinSize     int
	MaxPlaytime int
}

type CleanableFile struct {
	Path          string
	ModTime       time.Time
	Size          int64
	PathsToRemove []string
}

// Removable is to be implemented by any file
type Removable interface {
	// GetPaths returns the paths to the files that needs to be deleted to "fully" delete the given file.
	GetPaths() []string
}

func (f CleanableFile) Remove() error {
	return trash.Put(f.PathsToRemove...)
}

func (f CleanableFile) Exclude() {
	config.ExcludeFolder(f.Path)
}

func Clean(args Args) []CleanableFile {
	minAge := time.Now().AddDate(0, 0, -args.MinAge)
	clutterAnalyzer := clutter.NewAnalyzer(
		clutter.WithSizeFilter(args.MinSize),
		clutter.WithMinAgeFilter(minAge),
	)
	gamesAnalyzer := games.NewAnalyzer(
		games.WithMaxPlaytime(time.Duration(args.MaxPlaytime)*time.Hour),
		games.WithLastPlayedBefore(minAge),
	)
	mediaAnalyzer := media.NewAnalyzer()
	files := clutterAnalyzer.Analyze(args.Root)
	// TODO: Could probably get all analyzers on the same format.
	var cleanables []CleanableFile
	for _, file := range files {
		cleanables = append(cleanables, CleanableFile{
			Path:          file.GetPath(),
			ModTime:       file.ModTime,
			Size:          file.Size,
			PathsToRemove: file.GetPaths(),
		})
	}
	gms, err := gamesAnalyzer.Analyze()
	if err == nil {
		for _, g := range gms {
			cleanables = append(cleanables, CleanableFile{
				Path:          g.Path,
				ModTime:       g.LastPlayed,
				Size:          g.Size,
				PathsToRemove: g.GetPaths(),
			})
		}
	}
	mediaFiles := mediaAnalyzer.Analyze(args.Root)
	for _, file := range mediaFiles {
		cleanables = append(cleanables, CleanableFile{
			Path:          file.GetPath(),
			ModTime:       file.ModTime,
			Size:          file.Size,
			PathsToRemove: file.GetPaths(),
		})
	}
	var filteredResult []CleanableFile
	for _, file := range cleanables {
		if !config.UserExcludedFolders[file.Path] {
			filteredResult = append(filteredResult, file)
		}
	}
	return filteredResult
}
