.PHONY: build test clean install run-demo help

# 变量定义
BINARY_NAME=sql-diff
BUILD_DIR=bin
MAIN_PATH=cmd/sql-diff/main.go

# 默认目标
.DEFAULT_GOAL := help

## help: 显示帮助信息
help:
	@echo "可用的 make 命令："
	@echo "  make build      - 编译项目"
	@echo "  make test       - 运行所有测试"
	@echo "  make clean      - 清理构建产物"
	@echo "  make install    - 安装到 GOPATH/bin"
	@echo "  make run-demo   - 运行演示脚本"
	@echo "  make all        - 编译、测试并安装"

## build: 编译项目
build:
	@echo "正在编译 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ 编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

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
