package cmd

import (
	"fmt"

	"github.com/shirou/gopsutil/disk"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A little test command",
	Long:  "A little test command",
	Run: func(cmd *cobra.Command, args []string) {
		diskStats, err := disk.Partitions(true)
		if err != nil {
			log.WithError(err).Error("Could not collect disk stats!")
		}

		u := make(map[string]*disk.UsageStat)

		for _, v := range diskStats {
			u[v.Mountpoint], _ = disk.Usage(v.Mountpoint)
		}

		fmt.Println(u)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
