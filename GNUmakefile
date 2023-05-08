TEST                ?= ./internal/provider/... ./internal/client/... ./internal/config/... ./internal/utils/...
PKG_NAME            ?= internal
GO_VER              ?= go
TEST_COUNT          ?= 1
ACCTEST_PARALLELISM ?= 20
ACCTEST_TIMEOUT     ?= 180m

default: build

build: fmtcheck
	$(GO_VER) install

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@$(GO_VER) mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w -l ./$(PKG_NAME) tools.go main.go

fmtcheck:
	@sh -c "'$(CURDIR)/.ci/scripts/gofmtcheck.sh'"

test: fmtcheck
	$(GO_VER) test $(TEST) -v $(TESTARGS) -timeout=5m

testacc: fmtcheck
	TF_ACC=1 $(GO_VER) test ./${PKG_NAME}/acceptance/... -v -count $(TEST_COUNT) -parallel $(ACCTEST_PARALLELISM) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT)

vet:
	@echo "go vet ."
	@go vet $$(go list ./...) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: \
	build \
	depscheck \
	fmt \
	fmtcheck \
	test \
	testacc \
	vet
