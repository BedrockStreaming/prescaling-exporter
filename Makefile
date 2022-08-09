# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: generate
generate: pkg/apis/prescaling.bedrock.tech/v1/types.go vendor/modules.txt ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	hack/update-codegen.sh

.PHONY: fmt
fmt: ## Run go fmt against code.
	goimports -local $$(go list -m) -w $$(find . -type f -iname '*.go' ! -path './vendor/*' -print)

.PHONY: vet
vet: ## Run go vet against code.
	go vet $$(go list -e -compiled -test=false -export=false -deps=false -find=false -tags= -- ./...)

.PHONY: test
test: generate fmt vet ## Run tests.
	go test -coverprofile cover.out ./...

##@ Build

.PHONY: build
build: generate fmt vet ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

##@ Build Dependencies

vendor/modules.txt: go.mod go.sum
	go mod vendor
