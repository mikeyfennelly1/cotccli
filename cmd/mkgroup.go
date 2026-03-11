package cmd

import (
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	groupName string
)

var mkgroup = &cobra.Command{
	Use:   "mkgroup",
	Short: "Create a new group.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		c := client.NewGroupControllerClient(cfg.GetWebAppBaseUrl())
		err = c.CreateGroup(groupName)
		if err != nil {
			log.Fatalf("failed to create group: %v\n", err)
		}
	},
}

func init() {
	mkgroup.Flags().StringVarP(&groupName, "name", "n", "", "Name of the new group")

	mkgroup.MarkFlagRequired("name")

	rootCmd.AddCommand(mkgroup)
}
