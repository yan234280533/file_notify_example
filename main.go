package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"runtime"
)

const TimeFormat = "2006-01-02 15:04:05"

func main() {
	var seconds int
	var enableSysrq bool // 新增开关变量
	pflag.IntVarP(&seconds, "seconds", "s", 10, "Number of seconds to sleep")
	pflag.BoolVarP(&enableSysrq, "sysrq", "r", false, "Enable sysrq-trigger on timeout") // 新增参数
	pflag.Parse()

	if seconds <= 0 {
		fmt.Println("Invalid time parameter, must be greater than 0")
		pflag.Usage()
		return
	}

	for {
		fmt.Printf("New process started 001: %v\n", time.Now().Format(TimeFormat))
		cmd := exec.Command("/bin/sleep", fmt.Sprintf("%d", seconds))
		startTime := time.Now()
		fmt.Printf("Executing %d-second command: %v\n", seconds, startTime.Format(TimeFormat))

		if err := cmd.Start(); err != nil {
			fmt.Printf("Command failed to start: %v\n", err)
			continue
		}

		// Get and log process PID
		pid := cmd.Process.Pid
		fmt.Printf("Process PID %d started\n", pid)

		done := make(chan error)
		go func() { done <- cmd.Wait() }()

	checkLoop:
		for {
			select {
			case err := <-done:
				if err != nil {
					fmt.Printf("Command execution error: %v\n", err)
				} else {
					endTime := time.Now()
					fmt.Printf("Command completed at: %v (duration: %v)\n",
						endTime.Format(TimeFormat),
						endTime.Sub(startTime).Round(time.Second))
				}
				break checkLoop
			case <-time.After(1 * time.Second):
				elapsed := time.Since(startTime)
				fmt.Printf("Command has been running: %v seconds\n", elapsed.Round(time.Second))
				// Timeout detection
				if elapsed > (time.Duration(seconds*2) * time.Second) {
					fmt.Printf("ERROR: Process %d timeout (expected %ds, elapsed %ds)\n",
						pid,
						seconds*2,
						int(elapsed.Seconds()))

					// 添加条件判断
					if enableSysrq && runtime.GOOS == "linux" {
						if err := os.WriteFile("/proc/sysrq-trigger", []byte("c"), 0200); err != nil {
							fmt.Printf("Failed to trigger sysrq: %v (try running with sudo)\n", err)
						} else {
							fmt.Println("Triggered kernel stack trace via sysrq")
						}
					}

					// New stack reading logic
					fmt.Printf("Process stack trace:\n%s\n", readProcStack(pid))
					fmt.Printf("Thread stacks:\n%s\n", readThreadStacks(pid))
				}
			}
		}
	}
}

// Add these new functions at the bottom of the file
func readProcStack(pid int) string {
	stackPath := fmt.Sprintf("/proc/%d/stack", pid)
	content, err := os.ReadFile(stackPath)
	if err != nil {
		return fmt.Sprintf("Failed to read process stack: %v", err)
	}
	return string(content)
}

func readThreadStacks(pid int) string {
	var builder strings.Builder
	taskDir := fmt.Sprintf("/proc/%d/task", pid)

	entries, err := os.ReadDir(taskDir)
	if err != nil {
		return fmt.Sprintf("Failed to read task directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			tid := entry.Name()
			stackPath := filepath.Join(taskDir, tid, "stack")
			content, err := os.ReadFile(stackPath)
			if err != nil {
				builder.WriteString(fmt.Sprintf("Thread %s stack error: %v\n", tid, err))
				continue
			}
			builder.WriteString(fmt.Sprintf("Thread %s stack:\n%s\n", tid, content))
		}
	}
	return builder.String()
}
