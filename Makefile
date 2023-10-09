ARCHITECTURES=386 amd64
LDFLAGS=-ldflags="-s -w"

default: build

all: vet test build build-view

vet:
	go vet

test:
	go test ./... -v -timeout 10m

build:
	$(foreach GOARCH,$(ARCHITECTURES),\
		$(shell export GOARCH=$(GOARCH))\
		$(shell go build $(LDFLAGS) -o ./bin/$(GOARCH)/webServer ./cmd/webServer/...)\
		$(shell go build $(LDFLAGS) -o ./bin/$(GOARCH)/dataServer ./cmd/dataServer/...)\
		$(shell go build $(LDFLAGS) -o ./bin/$(GOARCH)/util ./cmd/util/...)\
	)\

build-web:
	yarn install && yarn prod

run:
	cd bin && ./webServer

dev:
	cd bin && ./webServer -m development

dev-web:
	yarn install && yarn dev