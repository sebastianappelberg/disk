package storage

import (
	"github.com/sebastianappelberg/disk/pkg/util"
	"os"
	"sync"
	"time"
)

const (
	semaphoreLimit = 100
	fileChLimit    = 250
)

type File struct {
	Base    string
	Name    string
	Size    int64
	IsDir   bool
	ModTime time.Time
}

func (f File) GetPaths() []string {
	return []string{f.GetPath()}
}

func (f File) GetPath() string {
	return util.SimpleJoin(f.Base, f.Name)
}

type FilterDecision int

const (
	// Include tells the walker to include the current entry in the result.
	Include FilterDecision = 1 << iota
	// Continue tells the walker to proceed as normal to the next entry.
	Continue
	// Skip tells the walker to proceed to the next entry instead of recursively digging into the file.
	Skip
	// ShortCircuit tells the walker to stop completely.
	ShortCircuit
)

func (d FilterDecision) Includes(flag FilterDecision) bool {
	return d&flag != 0
}

type Filter func(file File) FilterDecision

func IdentityFilter(_ File) FilterDecision {
	return Include
}

type Mapper[T any] func(file File, siblings []os.DirEntry) T

func IdentityMapper(file File, _ []os.DirEntry) File {
	return file
}

type FileWalker[T any] struct {
	filter    Filter
	mapper    Mapper[T]
	wg        sync.WaitGroup
	semaphore chan struct{}
}

type FileWalkerOption[T any] func(*FileWalker[T])

func WithMapper[T any](mapper Mapper[T]) FileWalkerOption[T] {
	return func(a *FileWalker[T]) {
		a.mapper = mapper
	}
}

func WithDecisionFilter[T any](filter Filter) FileWalkerOption[T] {
	return func(a *FileWalker[T]) {
		a.filter = filter
	}
}

func NewFileWalker[T any](options ...FileWalkerOption[T]) *FileWalker[T] {
	a := &FileWalker[T]{
		filter:    IdentityFilter,
		wg:        sync.WaitGroup{},
		semaphore: make(chan struct{}, semaphoreLimit),
	}
	for _, option := range options {
		option(a)
	}
	return a
}

func (w *FileWalker[T]) GetFiles(root string) <-chan T {
	ch := make(chan T, fileChLimit)

	w.wg.Add(1)
	go w.getFiles(root, ch)

	go func() {
		w.wg.Wait()
		close(ch)
	}()

	return ch
}

// getFiles scans the directory specified by dir. Subdirectories are walked recursively in separate goroutines.
func (w *FileWalker[T]) getFiles(dir string, entryCh chan<- T) {
	defer w.wg.Done()
	w.semaphore <- struct{}{}
	defer func() { <-w.semaphore }()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		file := File{
			Base:    dir,
			Name:    e.Name(),
			IsDir:   e.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}

		decision := w.filter(file)
		if decision.Includes(Include) {
			entryCh <- w.mapper(file, entries)
		}

		if decision.Includes(ShortCircuit) {
			break
		}

		if !file.IsDir || decision.Includes(Skip) {
			// We don't need to dig deeper once we've gotten the skip decision or the entry isn't a folder.
			continue
		}

		w.wg.Add(1)
		go w.getFiles(file.GetPath(), entryCh)
	}
}
