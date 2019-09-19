package collector

import (
	"github.com/shirou/gopsutil/disk"
	log "github.com/sirupsen/logrus"
)

// DiskStats represents the disk metrics
type DiskStats struct {
	Partitions []disk.PartitionStat
	Usages     map[string]*disk.UsageStat
}

// Update collects the disk metrics
func (d *DiskStats) Update() {
	diskStats, err := disk.Partitions(true)
	if err != nil {
		log.WithError(err).Error("Could not collect disk stats!")
	}

	d.Partitions = diskStats
	u := make(map[string]*disk.UsageStat)

	for _, v := range diskStats {
		u[v.Mountpoint], _ = disk.Usage(v.Mountpoint)
	}
	d.Usages = u
}
