package util

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const pathSeparator = string(os.PathSeparator)

func Time(label string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s: %s\n", label, time.Since(start))
	}
}

// SimpleJoin takes only two arguments and joins them by
// os.PathSeparator. It's meant to be a more performant but less flexible version
// of filepath.Join when you know how dir and name looks like.
func SimpleJoin(dir, name string) string {
	return dir + pathSeparator + name
}

// GetDirectoryDepth calculates how deep `currentPath` is relative to `rootPath`.
func GetDirectoryDepth(rootPath, currentPath string) int {
	if rootPath == currentPath {
		return 0
	}
	// A single path separator means depth is 0.
	rootDepth := 0
	if rootPath != "\\" && rootPath != "/" {
		rootDepth = strings.Count(rootPath, "\\") + strings.Count(rootPath, "/")
	}
	currentDepth := strings.Count(currentPath, "\\") + strings.Count(currentPath, "/")
	return currentDepth - rootDepth
}

// mapCommonPath finds the longest shared directory for each path.
func mapCommonPath(paths []string) map[string]string {
	if len(paths) == 0 {
		return nil
	}

	// Step 1: Sort paths lexicographically to group similar paths together
	sort.Strings(paths)

	// Step 2: Find the global common parent directories
	commonRoots := findCommonRoots(paths)

	// Step 3: Assign each path to its determined common root
	commonPaths := make(map[string]string)
	for _, path := range paths {
		commonPaths[path] = findMatchingRoot(path, commonRoots)
	}

	return commonPaths
}

// findCommonRoots finds the deepest common root paths for groups of similar directories.
func findCommonRoots(paths []string) []string {
	var roots []string

	for _, path := range paths {
		matched := false

		// Check if this path belongs under an existing root
		for i, root := range roots {
			common := commonBase(root, path)
			if common == root {
				matched = true
				break
			} else if common != "" {
				// Replace with a more general common root
				roots[i] = common
				matched = true
				break
			}
		}

		// If no match found, add it as a new root
		if !matched {
			roots = append(roots, path)
		}
	}

	return roots
}

// findMatchingRoot finds the deepest common root for a given path.
func findMatchingRoot(path string, roots []string) string {
	longestMatch := ""

	for _, root := range roots {
		if strings.HasPrefix(path, root) {
			if len(root) > len(longestMatch) {
				longestMatch = root
			}
		}
	}

	return longestMatch
}

// commonBase finds the longest shared directory prefix between two paths.
func commonBase(path1, path2 string) string {
	sep := "\\"
	segments1 := strings.Split(path1, sep)
	segments2 := strings.Split(path2, sep)

	var commonSegments []string
	for i := 0; i < len(segments1) && i < len(segments2); i++ {
		if segments1[i] == segments2[i] {
			commonSegments = append(commonSegments, segments1[i])
		} else {
			break
		}
	}

	return strings.Join(commonSegments, sep)
}
