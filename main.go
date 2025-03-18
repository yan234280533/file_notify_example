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

		if err := cmd.Wait(); err != nil {
			fmt.Printf("命令执行错误: %v\n", err)
		} else {
			endTime := time.Now()
			fmt.Printf("命令完成于: %v (耗时: %v)\n",
				endTime.Format(TimeFormat),
				endTime.Sub(startTime).Round(time.Second))
		}
	}
}
