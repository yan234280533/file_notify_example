#!/bin/bash

# 镜像名称
IMAGE=$1

# 最大重试次数
MAX_RETRIES=10

# 当前重试次数
RETRY_COUNT=0

# 推送镜像的函数
push_image() {
    echo "推送镜像: $IMAGE (尝试第 $((RETRY_COUNT + 1)) 次)"
    docker push "$IMAGE"
    return $?
}

# 主逻辑
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    push_image
    if [ $? -eq 0 ]; then
        echo "镜像推送成功！"
        exit 0
    else
        echo "镜像推送失败，重试中..."
        RETRY_COUNT=$((RETRY_COUNT + 1))
        sleep 2 # 等待 2 秒后重试，避免频繁请求
    fi
done

echo "镜像推送失败，已达到最大重试次数 ($MAX_RETRIES 次)。"
exit 1