#!/bin/bash

set -x

# 无限循环
while true; do
    # 获取当前时间
    CURRENT_TIME=$(date "+%Y-%m-%d %H:%M:%S")

    # 将输出写入日志文件
        echo "=== $CURRENT_TIME ==="
        echo ""

    # 等待 10 秒
    #sleep 10
    read -t 10 -n 1 -s
done