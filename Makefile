GO_PACKAGES := $(shell go list -f '{{.Dir}}' ./...)

WORKING_DIR := $(shell pwd)

OUT_DIR := ${WORKING_DIR}/out

deps:
	go install github.com/onsi/ginkgo/v2/ginkgo
	go install go.uber.org/mock/mockgen@latest

dirs:
	mkdir -p ${OUT_DIR}

ginkgo-bootstrap:
	$(foreach package,													\
		$(GO_PACKAGES),													\
		cd ${package};													\
		test -f $(shell basename ${package} | tr '-' '_')_suite_test.go \
			|| ginkgo bootstrap;										\
		cd -;															\
	)

generate:
	go generate ./...

test: dirs
	ginkgo --cover -r --output-dir out/ -coverpkg=./...
	go tool cover -html ${OUT_DIR}/coverprofile.out -o ${OUT_DIR}/cover.html

