GINKGO_VERSION := v2.19.0
REPORTS_DIRECTORY := test/integration/reports
TAGS ?= ""
LOG_LEVEL := v

ifeq ($(CI), true)
	LOG_LEVEL = vv
endif

.PHONY: integration-tests-local
integration-tests-local: export MCP_SERVER_HOST=http://localhost:8080
integration-tests-local: integration-tests-nonprd

.PHONY: integration-tests-dev
integration-tests-dev: export ENVIRONMENT=dev
integration-tests-dev: export MCP_SERVER_HOST=https://mcp-server.template.euwe1.dev.cdo.system-monitor.com
integration-tests-dev: integration-tests-nonprd

.PHONY: integration-tests-stg
integration-tests-stg: export ENVIRONMENT=stg
integration-tests-stg: export API_SERVER_HOST=https://api-template.euwe1.stg.cdo.system-monitor.com
integration-tests-stg: integration-tests-nonprd

.PHONY: integration-tests-nonprd
integration-tests-nonprd: integration-tests-run

.PHONY: integration-tests-prd-apac
integration-tests-prd-apac: export API_SERVER_HOST=https://api-template.ap-southeast-2.prd.cdo.system-monitor.com
integration-tests-prd-apac: integration-tests-prd

# Germany
.PHONY: integration-tests-prd-euce
integration-tests-prd-euce: export API_SERVER_HOST=https://api-template.eu-central-1.prd.cdo.system-monitor.com
integration-tests-prd-euce: integration-tests-prd

# Remaining EU
.PHONY: integration-tests-prd-euwe
integration-tests-prd-euwe: export API_SERVER_HOST=https://api-template.eu-west-1.prd.cdo.system-monitor.com
integration-tests-prd-euwe: integration-tests-prd

.PHONY: integration-tests-prd-uswe
integration-tests-prd-uswe: export API_SERVER_HOST=https://api-template.us-west-2.prd.cdo.system-monitor.com
integration-tests-prd-uswe: integration-tests-prd

.PHONY: integration-tests-prd
integration-tests-prd: export TAGS=SMOKE
integration-tests-prd: export ENVIRONMENT=prd
integration-tests-prd: integration-tests-run

.PHONY: integration-tests-run
integration-tests-run:
	ginkgo -r -${LOG_LEVEL} -keep-going -compilers=1 --label-filter="${TAGS}" -procs=6 -randomize-all -timeout=15m -trace -junit-report=integration-tests-report.xml -json-report=integration-tests-report.json -output-dir=./${REPORTS_DIRECTORY} ./test/integration/...

.PHONY: integration-tests-install
integration-tests-install:
	go install github.com/onsi/ginkgo/v2/ginkgo@${GINKGO_VERSION}
