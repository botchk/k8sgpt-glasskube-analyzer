# Image URL to use all building/pushing image targets
IMG ?= docker.io/botchk/k8sgpt-glasskube-analyzer:latest

# Allows to override name and location of the "go" binary.
# This allows using different versions of Go, as outlined in the documentation: https://go.dev/doc/manage-install
GOCMD?=go

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell $(GOCMD) env GOBIN))
GOBIN=$(shell $(GOCMD) env GOPATH)/bin
else
GOBIN=$(shell $(GOCMD) env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

OUT_DIR ?= $(shell pwd)/dist

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: tidy fmt lint build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: tidy
tidy: ## Run go mod tidy
	$(GOCMD) mod tidy

.PHONY: fmt
fmt: ## Run go fmt against code
	$(GOCMD) fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	$(GOCMD) vet ./...

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

##@ Build

.PHONY: build
build: fmt lint ## Build binary
	$(GOCMD) build -o $(OUT_DIR)/ .

.PHONY: run
run: fmt ## Run from your host
	$(GOCMD) run .

# PLATFORMS defines the target platforms for the manager image be built to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - be able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enabled BuildKit. More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image to your registry (i.e. if you do not set a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To adequately provide solutions that are compatible with multiple platforms, you should consider using this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: ## Build and push docker image with cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name project-v3-builder
	$(CONTAINER_TOOL) buildx use project-v3-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm project-v3-builder
	rm Dockerfile.cross

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: deploy
deploy: helm ## Deploy to the K8s cluster specified in ~/.kube/config.
	$(HELM) install k8sgpt-glasskube-analyzer chart -n k8sgpt-glasskube-analyzer --create-namespace

.PHONY: undeploy
undeploy: helm ## Undeploy from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(HELM) uninstall k8sgpt-glasskube-analyzer -n k8sgpt-glasskube-analyzer

##@ Tools

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.60.1
golangci-lint: ## Download golangci-lint binary if it doesn't exist
	@[ -f $(GOLANGCI_LINT) ] || { \
	set -e ;\
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell dirname $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION) ;\
	}

HELM = $(shell pwd)/bin/helm-$(GOOS)-$(GOARCH)
HELM_VERSION ?= v3.11.3

helm: ## Download helm binary if it doesn't exist
	@[ -f $(HELM) ] || { \
	set -e ;\
	curl -L https://get.helm.sh/helm-$(HELM_VERSION)-$(GOOS)-$(GOARCH).tar.gz | tar xz; \
	mv $(GOOS)-$(GOARCH)/helm $(HELM); \
	chmod +x $(HELM); \
	rm -rf ./$(GOOS)-$(GOARCH)/; \
	}
