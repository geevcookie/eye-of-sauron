package collector

import (
	"github.com/shirou/gopsutil/cpu"
	log "github.com/sirupsen/logrus"
)

// CPUStats represents the CPU metrics
type CPUStats struct {
	Stats cpu.TimesStat
	Perc  cpu.TimesStat
}

// Update collects the CPU metrics and percentages
func (c *CPUStats) Update(s *Stats) {
	stats, err := cpu.Times(false)
	if err != nil {
		log.WithError(err).Error("Could not collect CPU stats!")
	}

	currentCPUStat := cpu.TimesStat{}
	for index, stat := range stats {
		if index == 0 {
			currentCPUStat = stat
		} else {
			currentCPUStat.Idle += stat.Idle
			currentCPUStat.System += stat.System
			currentCPUStat.User += stat.User
			currentCPUStat.Guest += stat.Guest
			currentCPUStat.GuestNice += stat.GuestNice
			currentCPUStat.Iowait += stat.Iowait
			currentCPUStat.Irq += stat.Irq
			currentCPUStat.Nice += stat.Nice
			currentCPUStat.Softirq += stat.Softirq
			currentCPUStat.Steal += stat.Steal
			currentCPUStat.Stolen += stat.Stolen
		}
	}

	currentCPUPercStat := cpu.TimesStat{}
	if c.Stats.System != 0.00 {
		currentCPUPercStat = GetCPUPercentage(c.Stats, currentCPUStat)
	}

	s.previousCPU = *c
	c.Stats = currentCPUStat
	c.Perc = currentCPUPercStat
}

// GetCPUPercentage calculates the percentages of the collected CPU times
func GetCPUPercentage(previous cpu.TimesStat, current cpu.TimesStat) cpu.TimesStat {
	allDelta := current.Total() - previous.Total()

	if allDelta == 0 {
		return current
	}

	calculate := func(field2 float64, field1 float64) float64 {
		return ((field1 - field2) / allDelta) * 100
	}

	cpuPerc := cpu.TimesStat{
		User:      calculate(previous.User, current.User),
		System:    calculate(previous.System, current.System),
		Idle:      calculate(previous.Idle, current.Idle),
		Nice:      calculate(previous.Nice, current.Nice),
		Iowait:    calculate(previous.Iowait, current.Iowait),
		Irq:       calculate(previous.Irq, current.Irq),
		Softirq:   calculate(previous.Softirq, current.Softirq),
		Steal:     calculate(previous.Steal, current.Steal),
		Guest:     calculate(previous.Guest, current.Guest),
		GuestNice: calculate(previous.GuestNice, current.GuestNice),
		Stolen:    calculate(previous.Stolen, current.Stolen),
	}

	return cpuPerc
}
