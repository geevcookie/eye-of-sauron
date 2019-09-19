package collector

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"eye-of-sauron/internal"
	log "github.com/sirupsen/logrus"
)

// Collector represents the metrics collector
type Collector struct {
	Stats     Stats
	Frequency int
}

// Stats represents the complete metrics object
type Stats struct {
	Memory      MemStats
	CPU         CPUStats
	Disk        DiskStats
	Load        LoadStats
	previousCPU CPUStats
}

type Metrics struct {
	CPU    CPUMetrics    `json:"cpu"`
	Memory MemoryMetrics `json:"memory"`
	Load   LoadMetrics   `json:"load"`
	Disks  []DiskMetrics `json:"disks"`
}

type CPUMetrics struct {
	System float64 `json:"system"`
	User   float64 `json:"user"`
	Idle   float64 `json:"idle"`
	Perc   float64 `json:"perc"`
}

type MemoryMetrics struct {
	Total string  `json:"total"`
	Used  string  `json:"used"`
	Free  string  `json:"free"`
	Perc  float64 `json:"perc"`
}

type DiskMetrics struct {
	Mount string  `json:"mount"`
	Used  string  `json:"used"`
	Free  string  `json:"free"`
	Perc  float64 `json:"perc"`
}

type LoadMetrics struct {
	Load1  float64 `json:"load_1"`
	Load5  float64 `json:"load_5"`
	Load15 float64 `json:"load_15"`
}

// NewCollector creates a new collector object from configuration
func NewCollector(cfg internal.Config) Collector {
	return Collector{
		Stats:     Stats{},
		Frequency: cfg.MetricsConfig.Frequency,
	}
}

// Start starts the collector on a configurable timer
func (c *Collector) Start(channel *chan Metrics) {
	tick := time.Tick(time.Duration(c.Frequency) * time.Second)

	// Run initial update to collect CPU times
	c.Stats.Update()

	for {
		select {
		case <-tick:
			c.Stats.Update()

			// Non blocking channel send
			select {
			case *channel <- c.Stats.Metrics():
				log.WithField("metrics", c.Stats.Metrics()).Info("Collected Metrics:")
			default:
				log.WithField("metrics", c.Stats.Metrics()).Info("Collected Metrics:")
			}
		}
	}
}

// Update runs the update method on each collector
func (s *Stats) Update() {
	s.CPU.Update(s)
	s.Memory.Update()
	s.Disk.Update()
	s.Load.Update()
}

func (s *Stats) Metrics() Metrics {
	m := Metrics{}

	// Disk Metrics
	m.Disks = []DiskMetrics{}
	for mount, metrics := range s.Disk.Usages {
		// Skip virtual mounts
		if metrics.Total == 0 {
			continue
		}

		d := DiskMetrics{
			Mount: mount,
			Used:  fmt.Sprintf("%0.2fGB", float64(metrics.Total-metrics.Free)/1024/1024/1024),
			Free:  fmt.Sprintf("%0.2fGB", float64(metrics.Free)/1024/1024/1024),
			Perc:  toFixed((float64(metrics.Total-metrics.Free) / float64(metrics.Total)) * 100, 2),
		}
		m.Disks = append(m.Disks, d)
	}

	// CPU Metrics
	m.CPU = CPUMetrics{
		System: toFixed(s.CPU.Perc.System, 2),
		User:   toFixed(s.CPU.Perc.User, 2),
		Idle:   toFixed(s.CPU.Perc.Idle, 2),
		Perc:   toFixed(s.CPU.Perc.Total()-s.CPU.Perc.Idle-s.CPU.Perc.Iowait, 2),
	}

	// Memory Metrics
	m.Memory = MemoryMetrics{
		Total: fmt.Sprintf("%0.2fGB", s.Memory.Total/1024/1024/1024),
		Used:  fmt.Sprintf("%0.2fGB", s.Memory.Used/1024/1024/1024),
		Free:  fmt.Sprintf("%0.2fGB", s.Memory.Free/1024/1024/1024),
		Perc:  toFixed((s.Memory.Used/s.Memory.Total)*100, 2),
	}

	// Load Metrics
	m.Load = LoadMetrics{
		Load1:  s.Load.AVG.Load1,
		Load5:  s.Load.AVG.Load5,
		Load15: s.Load.AVG.Load15,
	}

	return m
}

// String converts the collected metrics into a JSON object
func (s *Stats) String() string {
	tree := map[string]interface{}{
		"cpu": map[string]interface{}{
			"system": s.CPU.Perc.System,
			"user":   s.CPU.Perc.User,
			"idle":   s.CPU.Perc.Idle,
		},
		"memory": map[string]interface{}{
			"total": s.Memory.Total,
			"used":  s.Memory.Used,
			"free":  s.Memory.Free,
		},
		"disk": s.Disk.Usages,
	}

	response, _ := json.Marshal(tree)

	return string(response)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
