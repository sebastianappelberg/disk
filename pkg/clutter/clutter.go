package clutter

import (
	"github.com/sebastianappelberg/disk/pkg/config"
	"github.com/sebastianappelberg/disk/pkg/storage"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	defaultMinSize = 50 * storage.MegaByte
)

type AnalyzerOption func(*Analyzer)

func WithMinAgeFilter(minAge time.Time) AnalyzerOption {
	return func(a *Analyzer) {
		a.minAge = minAge
	}
}

func WithSizeFilter(size int) AnalyzerOption {
	return func(a *Analyzer) {
		if size >= 0 {
			a.minSize = int64(size) * storage.MegaByte
		}
	}
}

type Analyzer struct {
	walker         *storage.FileWalker[storage.File]
	sizeCalculator *storage.SizeCalculator
	minSize        int64
	minAge         time.Time
}

func NewAnalyzer(options ...AnalyzerOption) *Analyzer {
	defaultMinAge := time.Now().AddDate(0, 0, -90)
	sizeAnalyzer := &Analyzer{
		walker: storage.NewFileWalker[storage.File](
			storage.WithMapper(storage.IdentityMapper),
			storage.WithDecisionFilter[storage.File](decisionFilter),
		),
		sizeCalculator: storage.NewSizeCalculator(),
		minSize:        defaultMinSize,
		minAge:         defaultMinAge,
	}
	for _, option := range options {
		option(sizeAnalyzer)
	}
	return sizeAnalyzer
}

func decisionFilter(file storage.File) storage.FilterDecision {
	fileName := strings.ToLower(file.Name)
	if config.ClutterFolders[fileName] {
		return storage.Include | storage.Skip
	}
	if config.UnsafeFolders[fileName] {
		return storage.Skip
	}
	return storage.Continue
}

// Analyze returns a sorted list of candidates to delete.
func (a *Analyzer) Analyze(root string) []storage.File {
	foldersCh := a.walker.GetFiles(root)
	sizeCh := a.calculateFolderSizes(foldersCh)

	var files []storage.File
	for file := range sizeCh {
		if file.Size >= a.minSize && file.ModTime.Before(a.minAge) {
			files = append(files, file)
		}
	}
	// Path feels the most intuitive when reading through a list.
	sort.Slice(files, func(i, j int) bool {
		return files[i].GetPath() < files[j].GetPath()
	})
	return files
}

func (a *Analyzer) calculateFolderSizes(files <-chan storage.File) <-chan storage.File {
	var wg sync.WaitGroup
	ch := make(chan storage.File, 200)

	for file := range files {
		wg.Add(1)
		go func(f storage.File) {
			defer wg.Done()
			if f.IsDir {
				f.Size = a.sizeCalculator.GetSize(f.GetPath())
			}
			ch <- f
		}(file)
	}

	go func() {
		wg.Wait()
		a.sizeCalculator.Close()
		close(ch)
	}()
	return ch
}
