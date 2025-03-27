#!/bin/bash

set -x

# 无限循环
while true; do
    # 获取当前时间
    CURRENT_TIME=$(date "+%Y-%m-%d %H:%M:%S")

    # 执行 crictl ps 并将输出写入日志文件
        echo "=== $CURRENT_TIME ==="
        time crictl ps  2>&1
        echo ""

    # 等待 10 秒
    strace -tt -T sleep 10
done