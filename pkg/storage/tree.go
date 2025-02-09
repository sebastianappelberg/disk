package storage

import (
	"github.com/sebastianappelberg/disk/pkg/util"
	"os"
)

// Tree represents a folder and its subfolders.
type Tree struct {
	Name     string
	Children []Tree
	Size     int64
}

type TreeWalker struct {
	sizeCalculator *SizeCalculator
}

func NewTreeWalker() *TreeWalker {
	return &TreeWalker{
		sizeCalculator: NewSizeCalculator(),
	}
}

func (t *TreeWalker) GetTree(root string, maxDepth int) Tree {
	defer t.sizeCalculator.Close()
	rootNode := Tree{
		Name: root,
	}
	children, err := t.getChildren(root, maxDepth, 0)
	if err != nil {
		return rootNode
	}
	rootNode.Children = children
	for _, child := range rootNode.Children {
		rootNode.Size += child.Size
	}
	return rootNode
}

// getChildren recursively gets folders up to n levels deep.
func (t *TreeWalker) getChildren(currentDir string, maxDepth, currentDepth int) ([]Tree, error) {
	if currentDepth >= maxDepth {
		return nil, nil
	}

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return nil, err
	}

	var children []Tree

	for _, entry := range entries {
		if entry.IsDir() {
			subDir := util.SimpleJoin(currentDir, entry.Name())
			// Recursively get nestedChildren
			nestedChildren, err := t.getChildren(subDir, maxDepth, currentDepth+1)
			if err != nil {
				return nil, err
			}
			// Append the tree with its nestedChildren
			size := int64(0)
			if nestedChildren == nil {
				// If it's a leaf node then explicitly get the full child size.
				size = t.sizeCalculator.GetSize(subDir)
			}
			for _, child := range nestedChildren {
				size += child.Size
			}
			child := Tree{
				Name:     entry.Name(),
				Children: nestedChildren,
				Size:     size,
			}
			children = append(children, child)
		} else {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			children = append(children, Tree{
				Name: entry.Name(),
				Size: info.Size(),
			})
		}
	}

	return children, nil
}
