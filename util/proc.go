package util

import (
	"time"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	// god forbid i have no spacing between my imports goland
	"github.com/shirou/gopsutil/v4/cpu"
)

func MemoryUsage() (int64, error) {
	v, err := mem.VirtualMemory()

	if err != nil {
		return 0, err
	} else {
		return int64(v.Used), err
	}
}

func CpuUsage() (map[int]float64, float64, error) {
	c, e := cpu.Percent(time.Second, true)
	var totalUsage float64
	perCore := make(map[int]float64)
	for i, v := range c {
		totalUsage += v
		perCore[i] = v
	}
	return perCore, totalUsage / float64(len(c)), e
}

func DiskUsage() (int64, error) {
	v, err := disk.Usage("/")
	if err != nil {
		return 0, err
	} else {
		return int64(v.Used), err
	}
}
