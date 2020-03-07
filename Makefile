APPBIN          := $(shell ls ./cmd/)
COVERAGEOUT     := coverage.out
COVERAGETMP     := coverage.tmp
BINPATH         := $(CURDIR)/bin
GOPATH          := $(GOPATH)
GOPKGS          := $(shell go list ./... | grep -v /vendor/ 2>/dev/null)
VERSION         := $(shell git rev-parse --short=8 HEAD)
DISTBUCKETPATH  := ring-distributions/linked-event-activators
DOCKER_REGISTRY := docker.svc.ring.com

# This repo's root import path (under GOPATH).
PKG := github.com/EdisonJunior/linked-event-activators

GOTOOLDIR := $(shell go env GOTOOLDIR)
GODOC     := $(GOTOOLDIR)/godoc
LINT      := $(GOBIN)/lint

$(GODOC)  : ; @GOPATH=$(ORIGGOPATH) go get -v golang.org/x/tools/cmd/godoc
$(LINT)   : ; @GOPATH=$(ORIGGOPATH) go get -v github.com/golang/lint/golint

.PHONY: vet
vet: ; @for pkg in $(GOPKGS); do go vet $$pkg || exit 1; done

.PHONY: lint
lint: $(LINT) ; @for src in $(GOSOURCES); do GOPATH=$(ORIGGOPATH) golint $$src || exit 1; done

.PHONY: fmt
fmt: ; @for src in $(GOSOURCES); do GOPATH=$(ORIGGOPATH) go fmt $$src; done

.PHONY: unit
unit:
	@echo 'mode: set' > $(COVERAGEOUT)
	@exitcode=0; \
	for pkg in $(GOPKGS); do \
		go test -v -race -tags=unit -coverprofile=$(COVERAGETMP) -covermode=atomic $$pkg; \
		testexitcode=$$?; \
		if [ $$testexitcode -ne 0 ]; then \
			exitcode=$$testexitcode; \
		fi; \
		if [ -f $(COVERAGETMP) ]; then \
			grep -v 'mode: set' $(COVERAGETMP) >> $(COVERAGEOUT); \
			rm $(COVERAGETMP); \
		fi; \
	done; \
	exit $$exitcode

.PHONY: integration
integration:
	@go test -v -race -tags=integration ./...

.PHONY: test
test: unit

.PHONY: build
build:
	@mkdir -p $(BINPATH)
	@for pkg in $(APPBIN); do \
		srvname=`basename $$pkg`; \
		binname=stream-source-$$pkg; \
		buildout=$(BINPATH)/$$binname; \
		CGO_ENABLED=0 GOOS=linux go build -o $$buildout -ldflags "-s -X main.version=$(VERSION)" -a ./cmd/$$pkg; \
	done

.PHONY: dep
dep:
	@dep ensure -update

.PHONY: proto
proto:
	@protoc \
		-I. \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:$(GOPATH)/src \
		--grpc-gateway_out=logtostderr=true:$(GOPATH)/src \
		--proto_path=$(GOPATH)/src \
		--swagger_out=logtostderr=true:$(GOPATH)/src \
		$(GOPATH)/src/github.com/gleach/kill/interfaces/server/grpc/*.proto
