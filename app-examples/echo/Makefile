# golang1.17 or latest

# machine x
GO_OS_ARCH=$(shell go version | awk '{ print $$4 }')
HOST_ARCH=$(shell echo "${GO_OS_ARCH}" | awk -F '/' '{ print $$2 }')
HOST_OS=$(shell echo "${GO_OS_ARCH}" | awk -F '/' '{ print $$1 }')

# go build env
GO_ENV_BASE=CGO_ENABLED=0 GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB="sum.golang.google.cn"
GO_ENV=GOOS=$(HOST_OS) GOARCH=$(HOST_ARCH) ${GO_ENV_BASE}

# building injected params
version_file_path = cmd
project_version=$(shell git describe --abbrev=8 --tags || git rev-parse --short=8 HEAD)
go_version=$(shell go version | awk '{ print $$3 }')
git_commit=$(shell git rev-parse --short=8 HEAD)
build_time=$(shell date "+%Y-%m-%d %H:%M:%S %Z")

#pkg_commit_name=${project_version}_${git_commit}
TARGET_DIR=target
servers=client_tcp client_ws server

# flags
#issue1221: https://github.com/spf13/cobra/issues/1221
#pr39: https://github.com/mittwald/mittnite/pull/39
BUILD_FLAGS = -ldflags "\
-X 'github.com/txchat/im-util/protocol/$(version_file_path).projectVersion=$(project_version)' \
-X 'github.com/txchat/im-util/protocol/$(version_file_path).goVersion=$(go_version)' \
-X 'github.com/txchat/im-util/protocol/$(version_file_path).gitCommit=$(git_commit)' \
-X 'github.com/txchat/im-util/protocol/$(version_file_path).buildTime=$(build_time)' \
-X 'github.com/txchat/im-util/protocol/$(version_file_path).osArch=$(GO_OS_ARCH)' \
-X 'google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn'"

help: ## Display this help screen
	@printf "Help doc:\nUsage: make [command]\n"
	@printf "[command]\n"
	@grep -h -E '^([a-zA-Z_-]|\%)+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean build build_% test_%

build:  ## 编译本机系统和指令集的可执行文件
	for server in $(servers) ; do \
		$(GO_ENV) go build $(BUILD_FLAGS) -v -o ${TARGET_DIR}/$${server}_${HOST_OS}_${HOST_ARCH} $${server}/main.go; \
	done

build_%:  ## 编译目标机器的可执行文件（例如: make build_linux_amd64）
	TAR_OS=$(shell echo $* | awk -F'_' '{print $$1}'); \
	TAR_ARCH=$(shell echo $* | awk -F'_' '{print $$2}'); \
	for server in $(servers) ; do \
		GOOS=$${TAR_OS} GOARCH=$${TAR_ARCH} ${GO_ENV_BASE} go build $(BUILD_FLAGS) -v -o ${TARGET_DIR}/$${server}_$${TAR_OS}_$${TAR_ARCH} $${server}/main.go; \
	done

clean: ## Remove previous build
	@rm -rf ${TARGET_DIR}
	@go clean
