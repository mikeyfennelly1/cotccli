package cmd

import (
	"fmt"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	producerType string
	group        string
)

var startCmd = &cobra.Command{
	Use:   "mkproducer",
	Short: "Interfaces with the /api/producer API to create a new producer",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}

		producerClient := client.NewProducerClient(fmt.Sprintf(cfg.GetWebAppBaseUrl()))
		createdProducer, err := producerClient.CreateProducer(client.NewProducer{Name: producerName, Group: group})
		if err != nil || createdProducer == nil {
			log.Fatalf("failed to register producer: %v", err)
		}
	},
}

func init() {
	startCmd.Flags().StringVarP(&producerName, "producer-name", "n", "", "an identifiable name for the the producer")
	startCmd.Flags().StringVarP(&group, "group", "g", "", "the group name to send messages to on the collector")

	startCmd.MarkFlagRequired("producer-name")
	startCmd.MarkFlagRequired("group")

	rootCmd.AddCommand(startCmd)
}
