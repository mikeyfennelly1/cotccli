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
	producerName string
)

var produce = &cobra.Command{
	Use:   "produce",
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

		reportingClient := client.NewReportingClient(cfg.GetWebAppBaseUrl())
		producerMetadata, err := reportingClient.GetProducerByName(producerName)
		if err != nil {
			log.Fatalf("failed to find producer %q: %v", producerName, err)
		}

		producer := reader.ToProducer()
		log.Infof("starting scheduled producer")
		collectorClient := client.CollectorClient{BaseUrl: cfg.GetCollectorBaseUrl()}
		err = producer.StartScheduledProducer(&collectorClient, producerMetadata)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listeningPort), nil))
	},
}

func init() {
	produce.Flags().StringVarP(&producerName, "producer-name", "n", "", "Name of the registered producer to look up")
	produce.Flags().StringVarP(&producerType, "producer-type", "t", "", "The type of producer to start")

	produce.MarkFlagRequired("producer-name")
	produce.MarkFlagRequired("producer-name")

	rootCmd.AddCommand(produce)
}
