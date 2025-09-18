# SSH AI Server Makefile

# 变量定义
BINARY_NAME=sshai
MODULAR_BINARY=sshai
CMD_DIR=cmd
PKG_DIR=pkg
CONFIG_FILE=config.yaml

# 版本信息
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.9.19")
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION := $(shell go version | cut -d' ' -f3)

# 编译标志
LDFLAGS := -X 'sshai/pkg/version.Version=$(VERSION)' \
           -X 'sshai/pkg/version.GitCommit=$(GIT_COMMIT)' \
           -X 'sshai/pkg/version.BuildTime=$(BUILD_TIME)' \
           -X 'sshai/pkg/version.GoVersion=$(GO_VERSION)'

# 生产环境编译标志（去除调试信息，优化大小）
RELEASE_LDFLAGS := $(LDFLAGS) -s -w

# 默认目标
.PHONY: all
all: build

# 构建开发版本
.PHONY: build
build:
	@echo "Building development version..."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	go build -ldflags "$(LDFLAGS)" -o $(MODULAR_BINARY) $(CMD_DIR)/main.go

# 构建生产版本（优化大小，去除调试信息）
.PHONY: build-release
build-release:
	@echo "Building release version..."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	go build -ldflags "$(RELEASE_LDFLAGS)" -o $(MODULAR_BINARY) $(CMD_DIR)/main.go
	@echo "Release build complete (optimized, debug info stripped)"

# 运行模块化版本
.PHONY: run
run: build
	@echo "Starting SSH AI Server (modular)..."
	./$(MODULAR_BINARY)

# 测试
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

.PHONY: htop
htop:
	@echo "Monitoring $(BINARY_NAME)..."
	@htop -p $$(pgrep $(BINARY_NAME) || (echo "$(BINARY_NAME) not running!" && exit 1))


# 清理
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) $(MODULAR_BINARY)
	rm -f host_key.pem

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 检查代码
.PHONY: vet
vet:
	@echo "Vetting code..."
	go vet ./...

# 创建配置文件（如果不存在）
.PHONY: config
config:
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "Creating default config file..."; \
		cp config.yaml.example $(CONFIG_FILE) 2>/dev/null || \
		echo "Please create $(CONFIG_FILE) manually"; \
	else \
		echo "Config file already exists"; \
	fi

# 开发环境设置
.PHONY: dev-setup
dev-setup: deps config
	@echo "Development environment setup complete"

# 显示版本信息
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build development version with debug info"
	@echo "  build-release - Build optimized release version (no debug info)"
	@echo "  run           - Run development version"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  config        - Create default config"
	@echo "  dev-setup     - Setup development environment"
	@echo "  version       - Show version information"
	@echo "  help          - Show this help"
