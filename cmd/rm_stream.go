package cmd

import (
	"fmt"
	"os"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	"github.com/spf13/cobra"
)

var rmStreamName string

var rmStreamCmd = &cobra.Command{
	Use:   "rmstream",
	Short: "Delete a stream on the reporting API",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		c := client.NewConsumerClient(cfg.GetWebAppBaseUrl())
		err = c.DeleteStream(rmStreamName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete stream: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rmStreamCmd.Flags().StringVar(&rmStreamName, "name", "", "Name of the stream to delete (required)")
	rmStreamCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(rmStreamCmd)
}
