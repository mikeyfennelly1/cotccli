package cmd

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/sysinfo"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the application server",
	Run: func(cmd *cobra.Command, args []string) {
		port := "9090"
		log.Infof("Started Application listening on port %s", port)
		reading, err := sysinfo.GetReading()
		if err != nil {
			log.Fatal("an error occurred reading sysinfo", err)
		}
		readingJson, err := json.Marshal(reading)
		if err != nil {
			log.Errorf("error marshalling json to byte array: %v", err)
		}
		log.Infof("read data: %s", readingJson)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
