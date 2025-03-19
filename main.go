package main

import (
	"fmt"
	"os/exec"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

func main() {
	for {
		fmt.Printf("开始启动: %v\n", time.Now().Format(TimeFormat))
		cmd := exec.Command("/bin/sleep", "10")
		startTime := time.Now()
		fmt.Printf("开始执行命令: %v\n", startTime.Format(TimeFormat))

		if err := cmd.Start(); err != nil {
			fmt.Printf("命令启动失败: %v\n", err)
			continue
		}

		done := make(chan error)
		go func() { done <- cmd.Wait() }()

		// 异步状态检查循环
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
				fmt.Printf("命令已运行: %v秒\n", time.Since(startTime).Round(time.Second))
			}
		}
	}
}
