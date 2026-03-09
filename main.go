package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"time"
	"tnt/util"
)

func cpu() (prometheus.Gauge, *prometheus.GaugeVec) {
	usagePerCore, avg, _ := util.CpuUsage()

	var (
		total = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_usage_average",
			Help: "Average CPU usage",
		})

		perCore = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cpu_usage_core",
			Help: "CPU Usage per core",
		}, []string{"core"})
	)

	///
	total.Set(avg)
	for core, usage := range usagePerCore {
		perCore.With(prometheus.Labels{"core": fmt.Sprintf("%d", core)}).Set(usage)
	}
	return total, perCore
}

func ram() prometheus.Gauge {
	mem, _ := util.MemoryUsage()
	var (
		gauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "memory_usage",
			Help: "Memory usage",
		})
	)

	gauge.Set(float64(mem))
	return gauge
}

func disk() prometheus.Gauge {
	disk, _ := util.DiskUsage()
	var (
		gauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "memory_usage",
			Help: "Memory usage",
		})
	)

	gauge.Set(float64(disk))
	return gauge
}

func main() {

}
