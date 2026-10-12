# 实验 1 示例服务的构建入口。本地和 GitHub Actions 调用同一组目标，
# 这样“CI 里怎么构建”和“本地怎么构建”不会漂移。需要 GNU make、GNU tar、Go。
#
#   make fmt-check vet test                    quality gate（不修改任何文件）
#   make package VERSION=1.0.0 BUILD_NO=7      打包：dist/enroll-conflict_1.0.0_linux_amd64.tar.gz
#
# VERSION、BUILD_NO、RUN_ID、COMMIT 是追溯信息，由调用方传入；缺省值表示“没有经过流水线”。
# BUILD_NO 取 github.run_number，它在每个工作流里各自从 1 计数（ci 和 release 都有“第 1 次”）；
# RUN_ID 取 github.run_id，全局唯一，用来指回具体的那一次运行。

NAME     := enroll-conflict
VERSION  ?= dev
BUILD_NO ?= local
RUN_ID   ?= local
COMMIT   ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)
GOOS     ?= linux
# 默认 amd64，与 GitHub 托管的 ubuntu-24.04 runner（x64）一致；在 arm64 主机上要运行本地构建的二进制时，
# 用 make package GOARCH=arm64（包名里的 amd64 会变成 arm64）。
GOARCH   ?= amd64

LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildNo=$(BUILD_NO)
DIST    := dist
PKG     := $(NAME)_$(VERSION)_$(GOOS)_$(GOARCH)

.PHONY: fmt-check vet test build package clean

fmt-check:
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then echo "以下文件需要 gofmt："; echo "$$out"; exit 1; fi

vet:
	go vet ./...

test:
	go test -count=1 ./...

build:
	mkdir -p $(DIST)/$(PKG)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/$(PKG)/$(NAME) .

# BUILD_INFO 不写构建时间：同样的输入应得到同样的输出（见指导书分析问题 3）。
# tar 固定文件顺序、属主和修改时间，所以归档本身也可重复；其中对哈希影响最大的是
# --mtime（文件的修改时间每次构建都不同）。gzip -n 让 gzip 压缩文件时不写文件名和时间戳，
# 这里数据来自管道，本来就不写，保留它是为了换成“压缩文件”的写法时仍然可重复。
package: build
	printf 'name=%s\nversion=%s\nbuild_no=%s\nrun_id=%s\ncommit=%s\ngo=%s\n' \
		'$(NAME)' '$(VERSION)' '$(BUILD_NO)' '$(RUN_ID)' '$(COMMIT)' "$$(go env GOVERSION)" > $(DIST)/$(PKG)/BUILD_INFO
	tar --sort=name --owner=0 --group=0 --numeric-owner --mtime=@0 -C $(DIST) -cf - $(PKG) | gzip -n > $(DIST)/$(PKG).tar.gz
	cd $(DIST) && sha256sum $(PKG).tar.gz > $(PKG).tar.gz.sha256

clean:
	rm -rf $(DIST)
