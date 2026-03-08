package cmd

import (
	"fmt"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/spf13/cobra"
)

var healthBaseUrl string

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check health of reporting, collector, and consumer services",
	Run: func(cmd *cobra.Command, args []string) {
		checks := []struct {
			name string
			fn   func() error
		}{
			{"reporting", client.NewReportingClient(healthBaseUrl).Health},
			{"collector", (&client.CollectorClient{BaseUrl: healthBaseUrl}).Health},
			{"consumer", client.NewConsumerClient(healthBaseUrl).Health},
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
	healthCmd.Flags().StringVar(&healthBaseUrl, "url", "http://localhost:8080", "Base URL for all services")
	rootCmd.AddCommand(healthCmd)
}
