#!/bin/bash

# 启用调试模式
set -x

mkdir -p ./backup

# 获取当前时间
CURRENT_TIME=$(date "+%Y-%m-%d %H:%M:%S")

# 将输出写入日志文件
echo "=== $CURRENT_TIME ==="

# 定义要备份的文件列表
files=("crictl_ps_log.txt" "example-linux-amd64.txt" "sda_t.txt" "strace_monitor_log.txt" "date_monitor_log.txt" "pid_monitor_log.txt" "top_monitor_log.txt" "cristace_log.txt" "datestrace_monitor_log.txt" "dateread_monitor_log.txt")

# 遍历文件列表
for file in "${files[@]}"; do
  if [[ -f $file ]]; then
    # 创建备份文件名
    backup_file="./backup/${file}.backup"

    # 备份文件
    cp "$file" "$backup_file"
    echo "Backed up $file to $backup_file"

    # 清空原文件内容
    truncate -s 0 $file
    echo "Cleared contents of $file"
  else
    echo "File $file does not exist, skipping."
  fi
done

# 关闭调试模式
set +x