BINARY_NAME=example
GO_BUILD=go build -o bin/${BINARY_NAME}-${GOOS}-${GOARCH} # 修改输出路径
MAIN_FILE=main.go
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# VERSION is currently based on the last commit
GIT_VERSION ?= $(shell git describe --tags --always)
# REGISTRY is the container registry to push
# into. The default is to push to the staging
# registry, not production.
REGISTRY?=registry-dev.vestack.sbuxcf.net/bci

# IMAGE is the image name of descheduler in the remote registry
IMAGE:=$(REGISTRY)/example:$(GIT_VERSION)

.PHONY: all build mac linux image-linux push-image deps run help  # 添加help到PHONY列表

# 新增help目标（放在文件末尾）
help:
	@echo "Available targets:"
	@echo "  all     - Build binaries for all platforms (macOS and Linux)"
	@echo "  build   - Build for current platform"
	@echo "  mac     - Build macOS binary"
	@echo "  linux   - Build Linux binary"
	@echo "  image-linux   - Build Linux image"
	@echo "  push-image   - Build push image"
	@echo "  clean   - Remove all binaries"
	@echo "  deps    - Install dependencies"
	@echo "  run     - Run the application with ARGS (e.g. make run ARGS=\"-s 5\")"
	@echo "  help    - Show this help message"

all: mac linux

build:
	@echo "Building for current platform..."
	@mkdir -p bin # 新增目录创建
	GOOS=${GOOS} GOARCH=${GOARCH} ${GO_BUILD} ${MAIN_FILE}

clean:
	@echo "Cleaning all binaries..."
	@rm -f bin/${BINARY_NAME}-* # 修改清理路径

mac:
	@echo "Building macOS binary..."
	@$(MAKE) build GOOS=darwin GOARCH=amd64

linux:
	@echo "Building Linux binary..."
	@$(MAKE) build GOOS=linux GOARCH=amd64

.PHONY: image-linux
image-linux: linux ## Build example image for linux
	docker build  --platform linux/amd64  -t $(IMAGE) . --build-arg https_proxy="https://goproxy.cn,direct" --build-arg APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=0 --build-arg DOCKER_TLS_VERIFY=0

deps:
	@echo "Installing dependencies..."
	go mod tidy

run:
	@echo "Starting application..."
	go run ${MAIN_FILE} ${ARGS}

.PHONY: push-image
push-image: image-linux ## Push example image
	./push_image.sh $(IMAGE)