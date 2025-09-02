SERVICE := keti3
MAIN := main.go

TAG := $(shell git tag -l --points-at HEAD)
AUTHOR := $(shell git log --pretty=format:"%an"|head -n 1)
VERSION := $(shell git rev-list HEAD | head -1)
BUILD_INFO := $(shell git log --pretty=format:"%s" | head -1)
BUILD_DATE := $(shell date +%Y-%m-%d\ %H:%M:%S)
CUR_PWD := $(shell pwd)

export GO111MODULE=on
export GOPROXY=https://goproxy.cn,http://go.xesv5.com/proxy,direct
export GOSUMDB=off

LD_FLAGS='-X "$(SERVICE)/version.TAG=$(TAG)" -X "$(SERVICE)/version.VERSION=$(VERSION)" -X "$(SERVICE)/version.AUTHOR=$(AUTHOR)" -X "$(SERVICE)/version.BUILD_INFO=$(BUILD_INFO)" -X "$(SERVICE)/version.BUILD_DATE=$(BUILD_DATE)"'

default: build

.PHONY: build
build:
	go build -v -ldflags $(LD_FLAGS) -gcflags "-N" -o ./bin/$(SERVICE) ./$(MAIN)

.PHONY: build-beishida
build-beishida:
	GOOS=linux GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -gcflags "-N" -o ./bin/$(SERVICE) ./$(MAIN)

.PHONY: air
air:
	@hash air > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin; \
	fi

	air init

.PHONY: auto-build
auto-build:
	air -d

.PHONY:clean
clean:
	rm bin/*

.PHONY:fmt
fmt:
	go fmt ./...

.PHONY:gen-local-cfg
gen-local-cfg:
	@echo "generate local config from development-config"
	@rm -rf $(CUR_PWD)/config/local/*
	@mkdir -p $(CUR_PWD)/config/local/
	@cp -R $(CUR_PWD)/config/development/* $(CUR_PWD)/config/local/
	@echo "generate success."
