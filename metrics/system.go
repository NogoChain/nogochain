package metrics

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// StartSystemMetrics 启动系统资源监控
func StartSystemMetrics() {
	go func() {
		for {
			// 收集CPU使用率
			if cpuPercent, err := cpu.Percent(0, false); err == nil && len(cpuPercent) > 0 {
				CPUUsage.Set(cpuPercent[0])
			}

			// 收集内存使用情况
			if memInfo, err := mem.VirtualMemory(); err == nil {
				MemoryUsage.Set(float64(memInfo.Used))
			}

			// 收集磁盘使用情况
			if diskInfo, err := disk.Usage("/"); err == nil {
				DiskUsage.Set(float64(diskInfo.Used))
			}

			// 每10秒收集一次
			time.Sleep(10 * time.Second)
		}
	}()
}
