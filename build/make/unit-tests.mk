.PHONY: go-test
go-test:
	@go test -v $(shell go list ./... | grep -v /test/)

.PHONY: go-test-with-report
go-test-with-report:
	@go test -v $(shell go list ./... | grep -v /test/) > test_output.txt; \
	e=$$?; \
	go-junit-report -in test_output.txt -set-exit-code -out report.xml; \
	exit $$e

.PHONY: jenkins-test-report
jenkins-test-report:
	go test -coverprofile=coverage.out  -v $(shell go list ./... | grep -v /test/) > test_output.txt; \
	e=$$?; \
	go-junit-report -in test_output.txt -set-exit-code -out report.xml; \
	gocover-cobertura < coverage.out > coverage.xml; \
	mkdir -p coverage-reports; \
	go tool cover -html=coverage.out -o coverage-reports/coverage.html; \
	exit $$e
