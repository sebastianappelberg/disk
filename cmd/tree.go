package cmd

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss/tree"
	"github.com/sebastianappelberg/disk/pkg/storage"
	"github.com/spf13/cobra"
)

func NewCmdTree() *cobra.Command {
	var depth int
	var sortBy string

	var cmd = &cobra.Command{
		Use:   "tree <path>",
		Short: "Print folders and files along with their sizes, in a tree structure.",
		Run: func(cmd *cobra.Command, args []string) {
			root := args[0]
			walker := storage.NewTreeWalker()
			folder := walker.GetTree(root, depth)
			t := buildTreeFromFolder(folder, sortBy)
			fmt.Println(t)
		},
	}

	cmd.Flags().IntVarP(&depth, "depth", "d", 1, "Depth of the tree structure.")
	cmd.Flags().StringVarP(&sortBy, "sort", "s", "name", "Sort by 'name' or 'size'.")

	return cmd
}

func sortChildren(children []storage.Tree, sortBy string) []storage.Tree {
	sorted := make([]storage.Tree, len(children))
	copy(sorted, children)

	if sortBy == "size" {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Size > sorted[j].Size // Sort descending by size
		})
	} else {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Name < sorted[j].Name // Sort ascending by name
		})
	}

	return sorted
}

func buildTreeFromFolder(folder storage.Tree, sortBy string) *tree.Tree {
	t := tree.Root(fmt.Sprintf("%s: %s", folder.Name, storage.FormatSize(folder.Size)))

	children := sortChildren(folder.Children, sortBy)

	for _, subfolder := range children {
		t.Child(buildTreeFromFolder(subfolder, sortBy))
	}
	return t
}
