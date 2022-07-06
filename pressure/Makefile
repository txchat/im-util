# golang1.17 or latest
TARGETDIR=target

# machine x
HOSTARCH=$(shell uname -m)
HOSTOS=$(shell uname -s | tr '[:upper:]' '[:lower:]')

GOBASE=CGO_ENABLED=0 GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB="sum.golang.google.cn"
GOENV=GOOS=$(HOSTOS) GOARCH=$(HOSTARCH) ${GOBASE}

# 编译时动态注入的参数
versionFilePath = cmd
projectVersion=$(shell git describe --abbrev=8 --tags || git rev-parse --short=8 HEAD)
goVersion=$(shell go version | awk '{ print $$3 }')
gitCommit=$(shell git rev-parse --short=8 HEAD)
buildTime=$(shell date "+%Y-%m-%d %H:%M:%S %Z")
osArch="$(HOSTOS)/$(HOSTARCH)"

pkgCommitName=${projectVersion}_${gitCommit}
TAR_MACHINE := ${HOSTOS}_${HOSTARCH}
APP_NAME := pressure

# flags
#issue1221: https://github.com/spf13/cobra/issues/1221
#pr39: https://github.com/mittwald/mittnite/pull/39
BUILD_FLAGS = -ldflags "\
-X 'github.com/txchat/im-util/protocol/$(versionFilePath).projectVersion=$(projectVersion)' \
-X 'github.com/txchat/im-util/protocol/$(versionFilePath).goVersion=$(goVersion)' \
-X 'github.com/txchat/im-util/protocol/$(versionFilePath).gitCommit=$(gitCommit)' \
-X 'github.com/txchat/im-util/protocol/$(versionFilePath).buildTime=$(buildTime)' \
-X 'github.com/txchat/im-util/protocol/$(versionFilePath).osArch=$(osArch)' \
-X 'google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn'"

help: ## Display this help screen
	@printf "Help doc:\nUsage: make [command]\n"
	@printf "[command]\n"
	@grep -h -E '^([a-zA-Z_-]|\%)+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean build build_% test_%

build:  ## 编译本机系统和指令集的可执行文件
	$(GOENV) go build $(BUILD_FLAGS) -v -o ${TARGETDIR}/${APP_NAME}_${TAR_MACHINE} main.go

build_%:  ## 编译目标机器的可执行文件（例如: make build_linux_amd64）
	TAR_OS=$(shell echo $* | awk -F'_' '{print $$1}'); \
	TAR_ARCH=$(shell echo $* | awk -F'_' '{print $$2}'); \
	GOOS=$${TAR_OS} GOARCH=$${TAR_ARCH} go build $(BUILD_FLAGS) -v -o ${TARGETDIR}/${APP_NAME}_$${TAR_OS}_$${TAR_ARCH} main.go

clean: ## Remove previous build
	@rm -rf ${TARGETDIR}
	@go clean
