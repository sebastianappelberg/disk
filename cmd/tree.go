package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/sebastianappelberg/disk/pkg/storage"
	"github.com/spf13/cobra"
)

func NewCmdTree() *cobra.Command {
	var depth int

	var cmd = &cobra.Command{
		Use:   "tree <path>",
		Short: "Print folders and files along with their sizes, in a tree structure.",
		Run: func(cmd *cobra.Command, args []string) {
			root := args[0]
			walker := storage.NewTreeWalker()
			folder := walker.GetTree(root, depth)
			t := buildTreeFromFolder(folder)
			fmt.Println(t)
		},
	}

	cmd.Flags().IntVarP(&depth, "depth", "d", 1, "Depth of the tree structure.")

	return cmd
}

func buildTreeFromFolder(folder storage.Tree) *tree.Tree {
	t := tree.Root(fmt.Sprintf("%s: %s", folder.Name, storage.FormatSize(folder.Size)))
	for _, subfolder := range folder.Children {
		t.Child(buildTreeFromFolder(subfolder))
	}
	return t
}
