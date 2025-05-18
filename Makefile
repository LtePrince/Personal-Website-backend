# 定义路径变量
CMD_DIR := ./cmd
MAIN_FILE := $(CMD_DIR)/main.go
TEST_DIR := ./test
CLIENT_MONITOR_FILE := $(TEST_DIR)/client_monitor.go

BUILD_DIR := ./build


# 运行目标程序
run:
	@go run $(MAIN_FILE)

# 运行客户端模拟程序
test:
	@go run $(CLIENT_MONITOR_FILE)

build:
	@go build -o $(BUILD_DIR)/myserver $(MAIN_FILE)

.PHONY: run test