package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/mem"

	"github.com/shirou/gopsutil/v4/cpu"
)

func main() {

	fmt.Println("Mem")
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(time.Second, true)
	var x float64 = 0
	var coreNum int = 0
	for i := 0; i < len(c); i++ {
		var usage float64 = c[i]
		x += usage
		coreNum++
		fmt.Printf("core %d - %f %% \n", i, usage)
	}
	fmt.Printf("Avg. usage: %f", x/float64(coreNum))
}
