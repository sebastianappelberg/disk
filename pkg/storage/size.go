package storage

import (
	"github.com/sebastianappelberg/disk/pkg/cache"
	"github.com/sebastianappelberg/disk/pkg/config"
	"os"
	"time"
)

type sizeCacheEntry struct {
	ModTime time.Time
	Size    int64
}

type SizeCalculator struct {
	walker *FileWalker[File]
	cache  *cache.Cache[sizeCacheEntry]
}

func NewSizeCalculator() *SizeCalculator {
	return &SizeCalculator{
		walker: NewFileWalker[File](
			WithDecisionFilter[File](IdentityFilter),
			WithMapper(IdentityMapper),
		),
		cache: cache.NewCache[sizeCacheEntry](config.GetAppDir(), "sizes"),
	}
}

func (s *SizeCalculator) GetSize(root string) int64 {
	fileInfo, err := os.Stat(root)
	if err != nil {
		return 0
	}
	size, ok := s.cache.Get(root)
	if ok && fileInfo.ModTime().Equal(size.ModTime) {
		return size.Size
	}
	fileCh := s.walker.GetFiles(root)

	var total int64
	for file := range fileCh {
		if !file.IsDir {
			total += file.Size
		}
	}
	s.cache.Put(root, sizeCacheEntry{ModTime: fileInfo.ModTime(), Size: total})
	return total
}

func (s *SizeCalculator) Close() {
	s.cache.Flush()
}
