.PHONY: lint
lint:
	golangci-lint run -c .golangci.yaml

.PHONY: lint-new
lint-new:
	golangci-lint run -c .golangci.yaml --new-from-rev=HEAD~1

.PHONY: lint-fix
lint-fix:
	golangci-lint run -c .golangci.yaml --fix

.PHONY: lint-fix-new
lint-fix-new:
	golangci-lint run -c .golangci.yaml --new-from-rev=HEAD~1 --fix
