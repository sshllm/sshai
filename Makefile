# SSH AI Server Makefile

# 变量定义
BINARY_NAME=sshai
MODULAR_BINARY=sshai
CMD_DIR=cmd
PKG_DIR=pkg
CONFIG_FILE=config.yaml

# 默认目标
.PHONY: all
all: build

# 构建模块化版本
.PHONY: build
build:
	@echo "Building modular version..."
	go build -o $(MODULAR_BINARY) $(CMD_DIR)/main.go

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

# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build modular version"
	@echo "  build-legacy - Build legacy version"
	@echo "  run          - Run modular version"
	@echo "  run-legacy   - Run legacy version"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  config       - Create default config"
	@echo "  dev-setup    - Setup development environment"
	@echo "  help         - Show this help"
