package cmd

import (
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	streamName string
)

var mkstream = &cobra.Command{
	Use:   "mkstream",
	Short: "Create a new stream.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		c := client.NewStreamControllerClient(cfg.GetWebAppBaseUrl())
		err = c.CreateStream(client.NewStream{
			Name:   streamName,
			Parent: "",
		})
		if err != nil {
			log.Fatalf("failed to create stream: %v\n", err)
		}
	},
}

func init() {
	mkstream.Flags().StringVar(&streamName, "name", "", "Name of the new stream")

	mkstream.MarkFlagRequired("name")

	rootCmd.AddCommand(mkstream)
}
