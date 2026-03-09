package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"os"
	"time"
	"tnt/util"
)

type CPU struct {
	Avg prometheus.Gauge
	Per *prometheus.GaugeVec
}

type mcStats struct {
	Online prometheus.Gauge
	Max    prometheus.Gauge
}

func cpu(c chan CPU) {
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
	c <- CPU{Avg: total, Per: perCore}
	return
}

func ram(c chan prometheus.Gauge) {
	mem, _ := util.MemoryUsage()
	var (
		gauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "memory_usage",
			Help: "Memory usage",
		})
	)

	gauge.Set(float64(mem))
	c <- gauge
	return
}

func disk(c chan prometheus.Gauge) {
	disk, _ := util.DiskUsage()
	var (
		gauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "disk_usage",
			Help: "disk usage",
		})
	)
	gauge.Set(float64(disk))
	c <- gauge
	return
}

func mc(c chan mcStats, address string, port int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	status, _ := util.Query(ctx, address, uint16(port))
	defer cancel()
	var (
		online = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "online",
			Help: "Players online",
		})

		maximum = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "maximum",
			Help: "Players maximum",
		})
	)

	maximum.Set(float64(status.Players.Max))
	online.Set(float64(status.Players.Online))
	c <- mcStats{Online: online, Max: maximum}
	return
}

func main() {
	var port = flag.Int("port", 25565, "MC Server Port")
	var address = flag.String("address", "localhost", "MC Server Address")
	var prom = flag.String("prom", "localhost", "Prometheus Address")
	flag.Parse()
	d := make(chan prometheus.Gauge)
	c := make(chan CPU)
	r := make(chan prometheus.Gauge)
	m := make(chan mcStats)
	go cpu(c)
	go ram(r)
	go disk(d)
	go mc(m, *address, *port)
	cpuData := <-c
	ramData := <-r
	diskData := <-d
	mcData := <-m
	err := push.New(*prom, "tnt").
		Collector(cpuData.Avg).
		Collector(cpuData.Per).
		Collector(ramData).
		Collector(diskData).
		Collector(mcData.Online).
		Collector(mcData.Max).
		Push()
	if err != nil {
		fmt.Println("failure", err)
		os.Exit(1)
	}
}
