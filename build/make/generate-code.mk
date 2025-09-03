# GENERATE MOCKS
.PHONY: generate-mocks
generate-mocks:
	@command -v mockery > /dev/null || (echo "https://github.com/vektra/mockery is required to generate mocks" && exit 1)
	@rm -fr mocks
	@mockery
	@mv mocks/internal mocks/internal_
