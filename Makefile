BUILD_FLAGS := -mod=vendor -v -ldflags "-w -s"
CODE_COVERAGE_DIR=./test

BUILD_DIR = ./bin
BINARY_NAME = iac-gen
BINARY_PATH = ${BUILD_DIR}/${BINARY_NAME}

default: cicd

# install deps
install:
	@which ginkgo > /dev/null || go install github.com/onsi/ginkgo/v2/ginkgo@latest
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.50.1

# clean
clean:
	@rm -rf ${BUILD_DIR} ${CODE_COVERAGE_DIR}

# build binary
build: clean
	@mkdir -p ${BUILD_DIR} > /dev/null
	CGO_ENABLED=1 go build ${BUILD_FLAGS} -o ${BINARY_PATH} ./cmd/iac-gen/main.go
	@echo "binary created at ${BINARY_PATH}"

# update go mod dependencies
dep:
	go mod tidy && go mod vendor

# run linter
lint:
	golangci-lint run -v --new-from-rev=origin/master

# create fakes for test
fakes:
	grep -rlw ./pkg -e 'go:generate' | xargs -n 1 -P 8 go generate

# run tests
test: clean test/unit

# run unit tests with ginkgo
test/unit:
	@mkdir -p ${CODE_COVERAGE_DIR} > /dev/null
	ginkgo -r -p --randomize-all --randomize-suites --fail-on-pending -mod=vendor --output-dir=${CODE_COVERAGE_DIR} --covermode=set --coverprofile=coverage.out

# replicate cicd 
cicd: install build test lint