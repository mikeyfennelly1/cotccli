package cmd

import (
	"context"
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
		port := "9090"
		log.Infof("Started Application listening on port %s", port)

		sysinfoChan := make(chan sysinfo.Message)
		go sysinfo.ScheduledProducer(context.Background(), sysinfoChan)

		go func() {
			for msg := range sysinfoChan {
				err := client.PushToAggregator(msg)
				if err != nil {
					log.Errorf("error writing to api: %v", err)
				}
			}
			log.Debugf("sysinfo channel closed, exiting worker")
		}()

		log.Fatal(http.ListenAndServe(":"+port, nil))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
