.PHONY: build test clean install run-demo help build-all build-linux build-darwin build-windows release

# 变量定义
BINARY_NAME=sql-diff
BUILD_DIR=bin
DIST_DIR=dist
MAIN_PATH=cmd/sql-diff/main.go
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go 编译参数
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"
GOFLAGS=-trimpath

# 平台和架构
PLATFORMS=linux darwin windows
ARCHS=amd64 arm64

# 默认目标
.DEFAULT_GOAL := help

## help: 显示帮助信息
help:
	@echo "SQL-Diff Makefile 命令列表"
	@echo ""
	@echo "开发命令:"
	@echo "  make build        - 编译当前平台的二进制文件"
	@echo "  make test         - 运行所有测试"
	@echo "  make clean        - 清理构建产物"
	@echo "  make install      - 安装到 GOPATH/bin"
	@echo "  make run-demo     - 运行演示脚本"
	@echo "  make fmt          - 格式化代码"
	@echo "  make lint         - 代码检查"
	@echo "  make deps         - 下载依赖"
	@echo ""
	@echo "跨平台编译:"
	@echo "  make build-all    - 编译所有平台 (Linux/macOS/Windows, amd64/arm64)"
	@echo "  make build-linux  - 编译 Linux (amd64 + arm64)"
	@echo "  make build-darwin - 编译 macOS (amd64 + arm64)"
	@echo "  make build-windows- 编译 Windows (amd64 + arm64)"
	@echo "  make release      - 编译并打包发布版本 (含校验和)"
	@echo ""
	@echo "其他:"
	@echo "  make all          - 编译、测试并安装"
	@echo "  make help         - 显示此帮助信息"
	@echo ""
	@echo "环境变量:"
	@echo "  VERSION           - 版本号 (默认: git tag)"
	@echo ""
	@echo "示例:"
	@echo "  make build                    # 编译当前平台"
	@echo "  make build-all                # 编译所有平台"
	@echo "  VERSION=v1.0.0 make release   # 打包 v1.0.0 发布版"

## build: 编译当前平台的二进制文件
build:
	@echo "正在编译 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ 编译完成: $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "✓ 版本: $(VERSION)"

## build-all: 编译所有平台和架构的二进制文件
build-all: clean
	@echo "开始编译所有平台..."
	@mkdir -p $(DIST_DIR)
	@$(MAKE) build-linux
	@$(MAKE) build-darwin
	@$(MAKE) build-windows
	@echo "✓ 所有平台编译完成"
	@ls -lh $(DIST_DIR)/

## build-linux: 编译 Linux 平台 (amd64 + arm64)
build-linux:
	@echo "编译 Linux 平台..."
	@GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=linux GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "✓ Linux 编译完成"

## build-darwin: 编译 macOS 平台 (amd64 + arm64)
build-darwin:
	@echo "编译 macOS 平台..."
	@GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "✓ macOS 编译完成"

## build-windows: 编译 Windows 平台 (amd64 + arm64)
build-windows:
	@echo "编译 Windows 平台..."
	@GOOS=windows GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@GOOS=windows GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PATH)
	@echo "✓ Windows 编译完成"

## release: 编译并打包发布版本 (包含校验和)
release: build-all
	@echo "生成 checksums..."
	@cd $(DIST_DIR) && shasum -a 256 * > checksums.txt
	@echo "✓ 发布包准备完成: $(DIST_DIR)/"
	@cat $(DIST_DIR)/checksums.txt

## test: 运行所有测试
test:
	@echo "正在运行测试..."
	@go test -v ./internal/parser
	@go test -v ./internal/differ
	@echo "✓ 所有测试通过"

## clean: 清理构建产物
clean:
	@echo "正在清理..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@go clean
	@echo "✓ 清理完成"

## install: 安装到 GOPATH/bin
install:
	@echo "正在安装 $(BINARY_NAME)..."
	@go install $(MAIN_PATH)
	@echo "✓ 安装完成"

## run-demo: 运行演示脚本
run-demo: build
	@echo "正在运行演示..."
	@chmod +x examples/demo.sh
	@./examples/demo.sh

## all: 编译、测试并安装
all: clean test build
	@echo "✓ 全部完成"

## deps: 下载依赖
deps:
	@echo "正在下载依赖..."
	@go mod download
	@go mod tidy
	@echo "✓ 依赖下载完成"

## fmt: 格式化代码
fmt:
	@echo "正在格式化代码..."
	@go fmt ./...
	@echo "✓ 代码格式化完成"

## lint: 代码检查
lint:
	@echo "正在进行代码检查..."
	@go vet ./...
	@echo "✓ 代码检查完成"
