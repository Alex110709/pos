# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: geth evm pixelzx all test lint fmt clean devtools help docker-build docker-build-hub compose-up compose-dev-up compose-logs compose-dev-logs

GOBIN = ./build/bin
GO ?= latest
GORUN = go run

#? geth: Build geth.
geth:
	$(GORUN) build/ci.go install ./cmd/geth
	@echo "Done building."
	@echo "Run \"$(GOBIN)/geth\" to launch geth."

#? evm: Build evm.
evm:
	$(GORUN) build/ci.go install ./cmd/evm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/evm\" to launch evm."

#? pixelzx: Build pixelzx.
pixelzx:
	$(GORUN) build/ci.go install ./cmd/pixelzx
	@echo "Done building."
	@echo "Run \"$(GOBIN)/pixelzx\" to launch pixelzx."

#? all: Build all packages and executables.
all:
	$(GORUN) build/ci.go install

#? test: Run the tests.
test: all
	$(GORUN) build/ci.go test

#? lint: Run certain pre-selected linters.
lint: ## Run linters.
	$(GORUN) build/ci.go lint

#? fmt: Ensure consistent code formatting.
fmt:
	gofmt -s -w $(shell find . -name "*.go")

#? clean: Clean go cache, built executables, and the auto generated folder.
clean:
	go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

#? docker-build: Build Docker image.
docker-build:
	docker build -t pixelzx-node .

#? docker-build-hub: Build Docker image for Docker Hub.
docker-build-hub:
	docker build -t yuchanshin/pixelzx-evm .

#? compose-up: Start production environment with Docker Compose.
compose-up:
	docker-compose up -d

#? compose-dev-up: Start development environment with Docker Compose.
compose-dev-up:
	docker-compose -f docker-compose.dev.yml up -d

#? compose-logs: Show production environment logs.
compose-logs:
	docker-compose logs -f

#? compose-dev-logs: Show development environment logs.
compose-dev-logs:
	docker-compose -f docker-compose.dev.yml logs -f

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

#? devtools: Install recommended developer tools.
devtools:
	env GOBIN= go install golang.org/x/tools/cmd/stringer@latest
	env GOBIN= go install github.com/fjl/gencodec@latest
	env GOBIN= go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	env GOBIN= go install ./cmd/abigen
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

#? help: Get more info on make commands.
help: Makefile
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@sed -n 's/^#?//p' $< | column -t -s ':' |  sort | sed -e 's/^/ /
