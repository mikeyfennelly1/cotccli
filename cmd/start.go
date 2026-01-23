package cmd

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/sysinfo"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the application server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Starting server on port %d", listeningPort)
		log.Infof("Collector agent: %s:%d",
			collectorAgentHostname,
			collectorAgentListeningPort,
		)

		sysinfoChan := make(chan sysinfo.Message)
		go sysinfo.ScheduledProducer(context.Background(), sysinfoChan)

		go func() {
			for msg := range sysinfoChan {
				err := client.PushToAggregator(msg, collectorAgentHostname, collectorAgentListeningPort)
				if err != nil {
					log.Errorf("error writing to api: %v", err)
				}
			}
			log.Debugf("sysinfo channel closed, exiting worker")
		}()

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listeningPort), nil))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
