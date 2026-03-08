package cmd

import (
	"fmt"
	"net/http"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/libproducer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	producerType    string
	producerSubname string
	stream          string
)

var startCmd = &cobra.Command{
	Use:   "producer",
	Short: "Starts a producer given a type, name and stream producerSubname.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}

		reader, err := libproducer.ReaderFactory(producerType, producerSubname)
		if err != nil {
			log.Fatalf("%v", err)
		}

		log.Debugf("Starting server of type %s on port %d", producerType, listeningPort)
		log.Debugf("Collector agent: %s:%d", collectorAgentHostname, collectorAgentListeningPort)

		consumerClient := client.NewConsumerClient(cfg.GetWebAppBaseUrl())
		createdProducer, err := consumerClient.CreateProducer(client.NewProducer{Name: reader.GetName(), StreamName: stream})
		if err != nil || createdProducer == nil {
			log.Fatalf("failed to register producer: %v", err)
		}

		producer := reader.ToProducer()
		log.Infof("starting scheduled producer ")
		collectorClient := client.CollectorClient{BaseUrl: cfg.GetCollectorBaseUrl()}
		producer.StartScheduledProducer(&collectorClient, createdProducer, stream)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listeningPort), nil))
	},
}

func init() {
	startCmd.Flags().StringVarP(&producerType, "type", "t", "", "Base URL of the reporting API")
	startCmd.Flags().StringVarP(&producerSubname, "producer-name", "n", "", "an identifiable name for the the producer")
	startCmd.Flags().StringVarP(&stream, "stream", "s", "", "the stream name to send messages to on the collector")

	startCmd.MarkFlagRequired("type")
	startCmd.MarkFlagRequired("producer-name")
	startCmd.MarkFlagRequired("stream")

	rootCmd.AddCommand(startCmd)
}
