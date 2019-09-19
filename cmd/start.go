package cmd

import (
	"eye-of-sauron/config"
	"eye-of-sauron/internal/metrics/collector"
	"eye-of-sauron/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Eye of Sauron server",
	Long:  "Start the Eye of Sauron server",
	Run: func(cmd *cobra.Command, args []string) {
		// Create communication channel
		channel := make(chan collector.Metrics)

		// Create metrics collector
		c := collector.NewCollector(cfg)
		go c.Start(&channel)

		// Create API Server
		s := server.NewServer()
		s.LoadRoutes(config.Routes(), &c, &channel)

		err := s.Run(cfg)
		if err != nil {
			log.WithError(err).Fatal("Could not start server!")
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
