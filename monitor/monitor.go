// 监控相关
package monitor

import (
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type Warn struct {
	Name    string // 什么告警
	Content string // 告警内容
}

// 超出 num 使用率持续 3 分钟（每 10s 采样一次 ） CPU 告警
func CpuUsage(num float64) (Warn, error) {
	var warn Warn
	var tfList [18]bool

	for i := range tfList {
		v, err := cpu.Percent(10*time.Millisecond, false)
		if err != nil {
			return warn, err
		}

		if v[0] > num {
			tfList[i] = true
		}
		time.Sleep(10 * time.Second)
	}

	for _, e := range tfList {
		if !e {
			return warn, nil
		}
	}
	warn = Warn{"cpu", "使用率超过 " + strconv.FormatFloat(num, 'g', 2, 64) + "% 持续一分钟"}
	return warn, nil
}

// 超出 num 使用率，内存告警
func NumUsage(num float64) (Warn, error) {
	var warn Warn
	v, err := mem.VirtualMemory()
	if err != nil {
		return warn, err
	}

	used := 100 - float64(v.Available)/float64(v.Total)*100
	if used > num {
		warn = Warn{"内存", "使用率超过 " + strconv.FormatFloat(num, 'g', 2, 64) + "%"}
		return warn, nil
	} else {
		return warn, nil
	}
}

// 超出 num 使用率，磁盘告警
func DiskUsage(num float64) ([]Warn, error) {
	var warnList []Warn
	partitions, err := disk.Partitions(false)
	if err != nil {
		return warnList, err
	}

	for _, e := range partitions {
		info, err := disk.Usage(e.Mountpoint)
		if err != nil {
			return warnList, err
		}

		used := float64(info.Used) / float64(info.Total) * 100
		if used > num {
			warnList = append(warnList, Warn{e.Device, "使用率超过 " + strconv.FormatFloat(num, 'g', 2, 64) + "%"})
		}
	}
	return warnList, nil
}
