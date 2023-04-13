TEST                ?= ./...
PKG_NAME   ?= internal
GO_VER     ?= go
TEST_COUNT ?= 1
ACCTEST_PARALLELISM ?= 20
ACCTEST_TIMEOUT     ?= 180m

default: build

build: fmtcheck
	$(GO_VER) install

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w -l ./$(PKG_NAME) tools.go main.go

fmtcheck:
	@sh -c "'$(CURDIR)/.ci/scripts/gofmtcheck.sh'"

test: fmtcheck
	$(GO_VER) test $(TEST) -v $(TESTARGS) -timeout=5m

testacc: fmtcheck
	TF_ACC=1 $(GO_VER) test ./${PKG_NAME}/... -v -count $(TEST_COUNT) -parallel $(ACCTEST_PARALLELISM) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT)

vet:
	@echo "go vet ."
	@go vet $$(go list ./...) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: build test testacc vet
