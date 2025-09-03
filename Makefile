project_root="$(CURDIR)"

# PRE-COMMIT
include build/make/pre-commit.mk

# LOCAL CLUSTER
include build/make/cluster.mk
include build/make/kubeconfig.mk

# QAPI / MOCKS
include build/make/generate-code.mk

# TERRAFORM
include build/make/terraform.mk

# TEST
include build/make/integration-tests.mk
include build/make/unit-tests.mk

# LINTING
include build/make/lint.mk

# DOCUMENTATION
adr-create:
	npx adr-tool create $(title)

# HOUSEKEEPING
.run-brew-upgrade:
	@brew update
	@brew bundle --file ./Brewfile

# BUILD
.PHONY: build
build-local:
	@kubectx $(KUBECONFIG_LOCAL)
	@skaffold build -f skaffold.yaml

setup-k8s-contexts:
	@./scripts/setup-k8s-contexts.sh

# Helper function to update local tooling within the repo
run-housekeeping: .run-brew-upgrade pre-commit-update

# LOCAL DEVELOPMENT
dev dev-auto-local:
	@kubectx $(KUBECONFIG_LOCAL)
# remove --cache-artifacts=false to use cached images, this is in to run the prebuild hooks.
	@skaffold dev -f skaffold.yaml --cache-artifacts=false --cleanup=true

.PHONY: version
version:
	@echo "${VERSION}"

.PHONY: go-env
go-env:
	@go env -w GOPRIVATE=github.com/nable-fusion/*,github.com/logicnow/*
	@go mod download

# TODO: include a way of updating this template after a new repo has been generated from it
