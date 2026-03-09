package cmd

import (
	"fmt"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var lsProducersStreamName string

var lsProducersCmd = &cobra.Command{
	Use:   "lsproducers",
	Short: "List all registered producers.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		reportingClient := client.NewReportingClient(cfg.GetWebAppBaseUrl())

		if lsProducersStreamName != "" {
			streamId, err := reportingClient.GetStreamUUIDByName(lsProducersStreamName)
			if err != nil {
				log.Fatalf("failed to find stream %q: %v", lsProducersStreamName, err)
			}

			producers, err := reportingClient.GetProducersByStreamId(streamId)
			if err != nil {
				log.Fatalf("failed to get producers for stream %q: %v", lsProducersStreamName, err)
			}

			for _, p := range producers {
				fmt.Printf("%s\t%s\n", p.UUID, p.ProducerName)
			}
			return
		}

		producers, err := reportingClient.GetProducers()
		if err != nil {
			log.Fatalf("failed to get producers: %v", err)
		}

		for _, p := range producers {
			fmt.Printf("%s\t%s\n", p.UUID, p.ProducerName)
		}
	},
}

func init() {
	lsProducersCmd.Flags().StringVarP(&lsProducersStreamName, "stream", "s", "", "Stream name to list producers for")
	rootCmd.AddCommand(lsProducersCmd)
}
