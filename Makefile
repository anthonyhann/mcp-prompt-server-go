.PHONY: build run clean test install deps

# 变量定义
BINARY_NAME=mcp-prompt-server
GO_FILES=$(shell find . -name "*.go" -type f)
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v2.0.0")
BUILD_TIME=$(shell date "+%Y-%m-%d %H:%M:%S")
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X 'main.BuildTime=$(BUILD_TIME)' -X main.CommitHash=$(COMMIT_HASH)"

# 默认目标
all: build

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# 构建二进制文件
build: deps
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go

# 构建用于生产的二进制文件（优化版本）
build-prod: deps
	@echo "Building production $(BINARY_NAME)..."
	CGO_ENABLED=0 go build -a -installsuffix cgo $(LDFLAGS) -o bin/$(BINARY_NAME) main.go

# 运行程序
run: build
	@echo "Starting MCP Prompt Server..."
	./bin/$(BINARY_NAME)

# 开发模式运行（实时重载）
dev:
	@echo "Starting development server with auto-reload..."
	go run main.go

# 测试
test:
	@echo "Running tests..."
	go test -v ./...

# 测试覆盖率
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# 检查代码格式
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 代码静态分析
vet:
	@echo "Running go vet..."
	go vet ./...

# 安装到系统
install: build-prod
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp bin/$(BINARY_NAME) /usr/local/bin/

# 卸载
uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# 创建发布包
package: build-prod
	@echo "Creating release package..."
	mkdir -p dist
	tar -czf dist/$(BINARY_NAME)-$(VERSION)-$(shell uname -s)-$(shell uname -m).tar.gz \
		bin/$(BINARY_NAME) \
		prompts/ \
		README.md \
		LICENSE

# 显示帮助信息
help:
	@echo "Available commands:"
	@echo "  build        - Build the binary"
	@echo "  build-prod   - Build optimized binary for production"
	@echo "  run          - Build and run the server"
	@echo "  dev          - Run in development mode"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Clean build files"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run static analysis"
	@echo "  install      - Install to system"
	@echo "  uninstall    - Remove from system"
	@echo "  package      - Create release package"
	@echo "  deps         - Install dependencies"
	@echo "  help         - Show this help" 