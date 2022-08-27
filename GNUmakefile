TEST            ?= $$(go list ./... | grep -v /sdk$$)
GOFMT_FILES     ?= $$(find . -name '*.go')
PKG_NAME        ?= sendgrid
ACCTEST_TIMEOUT ?= 10m
TEST_COUNT      ?= 1

ifneq ($(origin TESTS), undefined)
	RUNARGS = -run='$(TESTS)'
endif

ifneq ($(origin SWEEPERS), undefined)
	SWEEPARGS = -sweep-run='$(SWEEPERS)'
endif

default: build

build: fmtcheck
	go install
	$(MAKE) --directory=scripts doc

test: fmtcheck
	@go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	@if [ "$(TESTARGS)" = "-run=TestAccXXX" ]; then \
		echo ""; \
		echo "Error: Skipping example acceptance testing pattern. Update PKG and TESTS for the relevant *_test.go file."; \
		echo ""; \
		echo "For example if updating internal/service/acm/certificate.go, use the test names in internal/service/acm/certificate_test.go starting with TestAcc and up to the underscore:"; \
		echo "make testacc TESTS=TestAccACMCertificate_ PKG=acm"; \
		echo ""; \
		echo "See the contributing guide for more information: https://github.com/hashicorp/terraform-provider-aws/blob/main/docs/contributing/running-and-writing-acceptance-tests.md"; \
		exit 1; \
	fi
	TF_ACC=1 go test ./$(PKG_NAME)/... -v -count $(TEST_COUNT) $(RUNARGS) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT)

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)
	$(MAKE) --directory=scripts $@

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint: golangci-lint

golangci-lint:
	@echo "==> Checking source code against golangci-lint..."
	@golangci-lint run ./$(PKG_NAME)/...
	@$(MAKE) --directory=scripts $@

sweep:
	@rm -rf "$(CURDIR)/dist"
	@$(MAKE) --directory=scripts $@

test-release:
	goreleaser --snapshot --skip-publish --rm-dist

release:
	goreleaser release --rm-dist

.PHONY: build test testacc fmt fmtcheck lint golangci-lint sweep test-release release
