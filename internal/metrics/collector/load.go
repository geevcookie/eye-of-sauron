package collector

import (
	"github.com/shirou/gopsutil/load"
	log "github.com/sirupsen/logrus"
)

// LoadStat represents the load metrics
type LoadStats struct {
	AVG *load.AvgStat
}

// Update collects the load metrics
func (l *LoadStats) Update() {
	loadStats, err := load.Avg()
	if err != nil {
		log.WithError(err).Error("Could not collect load stats!")
	}

	l.AVG = loadStats
}