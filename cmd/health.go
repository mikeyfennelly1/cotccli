package cmd

import (
	"fmt"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check health of reporting, collector, and consumer services",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		checks := []struct {
			name string
			fn   func() error
		}{
			{"reporting", client.NewReportingClient(cfg.GetWebAppBaseUrl()).Health},
			{"collector", (&client.CollectorClient{BaseUrl: cfg.GetCollectorBaseUrl()}).Health},
			{"consumer", client.NewStreamControllerClient(cfg.GetWebAppBaseUrl()).Health},
		}

		allHealthy := true
		for _, svc := range checks {
			if err := svc.fn(); err != nil {
				fmt.Printf("[UNHEALTHY] %s: %v\n", svc.name, err)
				allHealthy = false
			} else {
				fmt.Printf("[OK]        %s\n", svc.name)
			}
		}

		if !allHealthy {
			fmt.Println("\nOne or more services are unhealthy.")
		}
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
