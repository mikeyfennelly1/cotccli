package cmd

import (
	"fmt"
	"os"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	"github.com/spf13/cobra"
)

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Print the stream hierarchy as a tree",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		c := client.NewReportingClient(cfg.GetWebAppBaseUrl())
		if err := c.GetStreamHierarchy(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to get stream hierarchy: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
}
