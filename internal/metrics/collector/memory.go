package collector

import (
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

// MemStats represents the memory metrics
type MemStats struct {
	Total float64
	Used  float64
	Free  float64
}

// Update collects the memory metrics
func (m *MemStats) Update() {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		log.WithError(err).Error("Could not collect memory stats!")
	}

	m.Total = float64(memStats.Total)
	m.Used  = float64(memStats.Used)
	m.Free  = float64(memStats.Total - memStats.Used)
}
