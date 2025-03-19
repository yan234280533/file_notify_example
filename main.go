package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/pflag"
)

const TimeFormat = "2006-01-02 15:04:05"

func main() {
	var seconds int
	pflag.IntVarP(&seconds, "seconds", "s", 10, "需要休眠的秒数")
	pflag.Parse()

	if seconds <= 0 {
		fmt.Println("无效的时间参数，必须大于0")
		pflag.Usage()
		return
	}

	for {
		fmt.Printf("开始启动: %v\n", time.Now().Format(TimeFormat))
		cmd := exec.Command("/bin/sleep", fmt.Sprintf("%d", seconds))
		startTime := time.Now()
		fmt.Printf("开始执行%d秒命令: %v\n", seconds, startTime.Format(TimeFormat))

		if err := cmd.Start(); err != nil {
			fmt.Printf("命令启动失败: %v\n", err)
			continue
		}

		// 获取并记录进程PID
		pid := cmd.Process.Pid
		fmt.Printf("进程PID %d 已启动\n", pid)

		done := make(chan error)
		go func() { done <- cmd.Wait() }()

		timeoutLogged := false
	checkLoop:
		for {
			select {
			case err := <-done:
				if err != nil {
					fmt.Printf("命令执行错误: %v\n", err)
				} else {
					endTime := time.Now()
					fmt.Printf("命令完成于: %v (耗时: %v)\n",
						endTime.Format(TimeFormat),
						endTime.Sub(startTime).Round(time.Second))
				}
				break checkLoop
			case <-time.After(1 * time.Second):
				elapsed := time.Since(startTime)
				fmt.Printf("命令已运行: %v秒\n", elapsed.Round(time.Second))
				// 超时检测逻辑
				if elapsed > time.Duration(seconds*2)*time.Second && !timeoutLogged {
					fmt.Printf("错误: 进程 %d 执行超时（预期%d秒，已运行%d秒）\n",
						pid,
						seconds*2,
						int(elapsed.Seconds()))
					timeoutLogged = true
				}
			}
		}
	}
}
