package cmd

import (
	"fmt"
	"github.com/ricochet2200/go-disk-usage/du"
	"github.com/sebastianappelberg/disk/pkg/storage"
	"github.com/spf13/cobra"
	"log"
	"os"
	"text/tabwriter"
)

func NewCmdUsage() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "usage",
		Short: "Print usage information for all available disks.",
		Run: func(cmd *cobra.Command, args []string) {
			disks, err := storage.GetAvailableDisks()
			if err != nil {
				log.Fatal(err)
			}
			totalUsed := uint64(0)
			totalSize := uint64(0)
			totalAvailable := uint64(0)
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
			fmt.Fprintln(w, "Disk\tSize\tUsed\tAvailable")
			for _, disk := range disks {
				diskUsage := du.NewDiskUsage(disk)
				totalUsed += diskUsage.Used()
				totalSize += diskUsage.Size()
				totalAvailable += diskUsage.Available()
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", disk, storage.FormatSize(diskUsage.Size()), storage.FormatSize(diskUsage.Used()), storage.FormatSize(diskUsage.Available()))
			}
			fmt.Fprintf(w, "Total:\t%s\t%s\t%s\n", storage.FormatSize(totalSize), storage.FormatSize(totalUsed), storage.FormatSize(totalAvailable))
			w.Flush()
		},
	}

	return cmd
}
