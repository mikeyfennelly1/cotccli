package cmd

import (
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/libproducer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	producerType string
	stream       string
)

var startCmd = &cobra.Command{
	Use:   "mkproducer",
	Short: "Starts a producer given a type, name and stream producerName.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}

		reader, err := libproducer.ReaderFactory(producerType, producerName)
		if err != nil {
			log.Fatalf("%v", err)
		}

		log.Debugf("Starting server of type %s on port %d", producerType, listeningPort)
		log.Debugf("Collector agent: %s:%d", collectorAgentHostname, collectorAgentListeningPort)

		consumerClient := client.NewStreamControllerClient(cfg.GetWebAppBaseUrl())
		createdProducer, err := consumerClient.CreateProducer(client.NewProducer{Name: reader.GetName(), StreamName: stream})
		if err != nil || createdProducer == nil {
			log.Fatalf("failed to register producer: %v", err)
		}
	},
}

func init() {
	startCmd.Flags().StringVarP(&producerType, "type", "t", "", "Base URL of the reporting API")
	startCmd.Flags().StringVarP(&producerName, "producer-name", "n", "", "an identifiable name for the the producer")
	startCmd.Flags().StringVarP(&stream, "stream", "s", "", "the stream name to send messages to on the collector")

	startCmd.MarkFlagRequired("type")
	startCmd.MarkFlagRequired("producer-name")
	startCmd.MarkFlagRequired("stream")

	rootCmd.AddCommand(startCmd)
}
