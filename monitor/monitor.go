// 监控相关
package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

// 重启停止容器
func RestartStopContainer() (string, error) {
	cmd := exec.Command("bash", "-c", "docker ps -a | grep Exited | awk  '{print $NF}'")
	stopContainer, err := cmd.CombinedOutput()
	if err != nil || len(stopContainer) == 0 {
		return "", err
	}

	outPut := ""
	for _, e := range strings.Split(strings.TrimSuffix(string(stopContainer), "\n"), "\n") {
		cmd = exec.Command("bash", "-c", fmt.Sprintf("docker restart %s", e))
		_, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		}
		outPut += fmt.Sprintf("%s docker restart %s\n", time.Now().Format("2006/01/02 15:04:05"), e)
	}
	return strings.TrimSuffix(outPut, "\n"), nil
}

type Process struct {
	FilePath string // 文件路径
	LogPath  string // 日志路径（缺省为 ${HOME}/monitorNohup.out）
}

// 重启停止进程
// FilePath 进程执行路径列表
func StartProcess(process *Process) (string, error) {
	cmd := exec.Command("bash", "-c", "pidof -q "+process.FilePath)
	_, err := cmd.CombinedOutput()
	// pidof 没找到返回状态 1，即错误，找到返回状态 0，err 为 nil
	if err == nil {
		return "", nil
	}

	// 日志路径缺省 ${HOME}/monitorNohup.out
	if process.LogPath == "" {
		process.LogPath = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "monitorNohup.out")
	}

	cmd = exec.Command("bash", "-c", fmt.Sprintf("nohup %s > %s 2>&1 &", process.FilePath, process.LogPath))
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s 启动 %s", time.Now().Format("2006/01/02 15:04:05"), process.FilePath), nil
}
