package system

import (
	"math"

	"github.com/shirou/gopsutil/cpu"
)

func CpuInfo() []cpu.InfoStat {
	info, _ := cpu.Info()
	return info
}

func CpuTimes(percpu bool) []cpu.TimesStat {
	info, _ := cpu.Times(percpu)
	return info
}

// 对比t1 CPU总时间与t2 Cpu总时间获取使用率
func CpuPercent(t1 cpu.TimesStat, t2 cpu.TimesStat) float64 {
	t1All, t1Busy := getAllBusy(t1)
	t2All, t2Busy := getAllBusy(t2)
	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 100
	}
	return math.Min(100, math.Max(0, (t2Busy-t1Busy)/(t2All-t1All)*100))
}

func getAllBusy(t cpu.TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal
	return busy + t.Idle, busy
}
