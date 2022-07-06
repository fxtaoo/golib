// 监控相关
package gomonitor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	if !warnTmp.Time.IsZero() && (warnTmp.Time.Sub(w.Time).Minutes() > notifyIntervalTime || w.Time.IsZero()) {
		*w = *warnTmp
		return true, nil
	}
	return false, nil
}

// 超出 num 使用率持续 1 分钟（每 5s 采样一次 ） CPU 告警
func CpuUsage(num float64) (*Warn, error) {
	var warn Warn
	var sampling [12]bool

	for range sampling {
		v, err := cpu.Percent(10*time.Millisecond, false)
		if err != nil {
			return &warn, err
		}

		// cpu 使用率小于 num，没有告警
		if v[0] < num {
			return &warn, nil
		}

		time.Sleep(5 * time.Second)
	}
	warn = Warn{time.Now(), fmt.Sprintf("cpu 使用率超过 %d%% 持续 1 分钟", int(num))}
	return &warn, nil
}

// 超出 num 使用率持续 1 分钟（每 5s 采样一次 ） 内存 告警
func NumUsage(num float64) (*Warn, error) {
	var warn Warn
	var sampling [12]bool

	for range sampling {
		v, err := mem.VirtualMemory()
		if err != nil {
			return &warn, err
		}

		used := 100 - float64(v.Available)/float64(v.Total)*100
		if used < num {
			return &warn, nil
		}

		time.Sleep(5 * time.Second)
	}

	warn = Warn{time.Now(), fmt.Sprintf("内存 使用率超过 %d%% 持续 1 分钟", int(num))}
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
			warn.Content += fmt.Sprintf("%s 使用率超过 %d%%  ", e.Device, int(num))
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
		outPut += fmt.Sprintf("docker restart %s\n", e)
	}
	return strings.TrimSuffix(outPut, "\n"), nil
}

// 检查进程是否运行
// 返回 []pid 字符串数组
func CheckProcess(filePath string) []string {
	cmd := exec.Command("bash", "-c", "pidof "+filePath)
	// pidof 没找到返回状态 1，即错误 err，找到返回状态 0，err 为 nil
	outPut, _ := cmd.CombinedOutput()
	if len(outPut) == 0 {
		return nil
	}
	return strings.Split(strings.TrimSuffix(string(outPut), "\n"), " ")
}

// 重启停止进程
// FilePath 执行文件绝对路径
func StartProcess(filePath, logPath string) (string, error) {
	if len(CheckProcess(filePath)) != 0 {
		return "", nil
	}

	// 日志路径缺省执行文件同目录下 monitorNohup.out
	if logPath == "" {
		logPath = filepath.Join(filepath.Dir(os.Args[0]), "monitorNohup.out")
	}

	if err := os.Chmod(filePath, 0755); err != nil {
		return "", err
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("nohup %s > %s 2>&1 &", filePath, logPath))
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("启动 %s", filePath), nil
}

// 进程不存在，执行脚本
func StartProcessScript(filePath, script string) string {
	if len(CheckProcess(filePath)) != 0 {
		return ""
	}
	cmd := exec.Command("bash", script)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("执行脚本 %s", script)
}
