package cmd

import (
	"fmt"
	"os"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	"github.com/spf13/cobra"
)

var rmgroup = &cobra.Command{
	Use:   "rmgroup",
	Short: "Delete a group on the reporting API",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		c := client.NewGroupControllerClient(cfg.GetWebAppBaseUrl())
		err = c.DeleteGroup(groupName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete group: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rmgroup.Flags().StringVarP(&groupName, "name", "n", "", "Name of the group to delete (required)")
	rmgroup.MarkFlagRequired("name")
	rootCmd.AddCommand(rmgroup)
}
