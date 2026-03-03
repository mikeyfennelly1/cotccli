package cmd

import (
	"fmt"
	"net/http"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/libproducer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the application server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Starting server on port %d", listeningPort)
		log.Infof("Collector agent: %s:%d", collectorAgentHostname, collectorAgentListeningPort)

		reader := libproducer.ReaderFactory("sysinfo", "1")
		producer := reader.ToProducer()
		producer.StartScheduledProducer()

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listeningPort), nil))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
