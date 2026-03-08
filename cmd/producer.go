package cmd

import (
	"fmt"
	"net/http"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/libproducer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	hostName     string
	port         int
	producerType string
	producerName string
	stream       string
)

var startCmd = &cobra.Command{
	Use:   "producer",
	Short: "Starts a producer given a type, name and stream producerName.",
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := libproducer.ReaderFactory(producerType, producerName)
		if err != nil {
			log.Fatalf("%v", err)
		}

		log.Debugf("Starting server of type %s on port %d", producerType, listeningPort)
		log.Debugf("Collector agent: %s:%d", collectorAgentHostname, collectorAgentListeningPort)

		producer := reader.ToProducer()
		log.Infof("starting scheduled producer ")
		collectorClient := client.NewCollectorClient(hostName, port)
		producer.StartScheduledProducer(collectorClient, stream)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listeningPort), nil))
	},
}

func init() {
	startCmd.Flags().StringVarP(&producerType, "type", "t", "", "Base URL of the reporting API")
	startCmd.Flags().StringVarP(&producerName, "producer-name", "n", "", "an identifiable name for the the producer")
	startCmd.Flags().StringVarP(&stream, "stream", "s", "", "the stream name to send messages to on the collector")
	startCmd.Flags().StringVarP(&hostName, "host", "H", "localhost", "the stream name to send messages to on the collector")
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "TCP port on the host that the collector API is reachable on")

	startCmd.MarkFlagRequired("type")
	startCmd.MarkFlagRequired("producer-name")
	startCmd.MarkFlagRequired("stream")

	rootCmd.AddCommand(startCmd)
}
