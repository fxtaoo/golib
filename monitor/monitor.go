// 监控相关
package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type Warn struct {
	Time    time.Time // 告警时间
	Content string    // 告警内容
}

// 检查告警间隔，大于将更新告警
// checkFun 检查函数，maxNum 最大使用率，notifyIntervalTime 告警间隔分钟
func (w *Warn) Check(checkFun func(float64) (*Warn, error), maxNum, notifyIntervalTime float64) (bool, error) {

	warnTmp, err := checkFun(maxNum)
	if err != nil {
		return false, err
	}

	// 两个 Warn 间隔时间是否超过 notifyIntervalTime 分钟
	if warnTmp.Time.Sub(w.Time).Minutes() > notifyIntervalTime {
		*w = *warnTmp
		return true, nil
	}
	return false, nil
}

// 超出 num 使用率持续 3 分钟（每 10s 采样一次 ） CPU 告警
func CpuUsage(num float64) (*Warn, error) {
	var warn Warn
	var sampling [18]bool

	for range sampling {
		v, err := cpu.Percent(10*time.Millisecond, false)
		if err != nil {
			return &warn, err
		}

		// cpu 使用率小于 num，没有告警
		if v[0] < num {
			return &warn, nil
		}
		time.Sleep(10 * time.Second)
	}
	warn = Warn{time.Now(), fmt.Sprintf("cpu 使用率超过 %d%% 持续 3 分钟\n", int(num))}
	return &warn, nil
}

// 超出 num 使用率持续 3 分钟（每 10s 采样一次 ） 内存 告警
func NumUsage(num float64) (*Warn, error) {
	var warn Warn
	var sampling [18]bool

	for range sampling {
		v, err := mem.VirtualMemory()
		if err != nil {
			return &warn, err
		}

		used := 100 - float64(v.Available)/float64(v.Total)*100
		if used < num {
			return &warn, nil
		}
	}

	warn = Warn{time.Now(), fmt.Sprintf("内存 使用率超过 %d%% 持续 3 分钟\n", int(num))}
	return &warn, nil
}

// 超出 num 使用率，磁盘 告警
func DiskUsage(num float64) (*Warn, error) {
	var warn Warn
	partitions, err := disk.Partitions(false)
	if err != nil {
		return &warn, err
	}

	for _, e := range partitions {
		info, err := disk.Usage(e.Mountpoint)
		if err != nil {
			return &warn, err
		}

		used := float64(info.Used) / float64(info.Total) * 100
		if used > num {
			warn.Content += fmt.Sprintf("%s 使用率超过 %d%% \n", e.Device, int(num))
			warn.Time = time.Now()
		}
	}
	return &warn, nil
}
