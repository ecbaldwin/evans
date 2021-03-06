SHELL := /bin/bash
VERSION := $(shell bump show meta/meta.go)

.PHONY: version
version:
	@echo "evans: $(VERSION)"


.PHONY: dep
dep:
ifeq ($(shell which dep 2>/dev/null),)
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

.PHONY: deps
deps: dep
	dep ensure

.PHONY: build
build: deps
	go build

.PHONY: test
test: vet lint unit-test e2e-test

.PHONY: unit-test
unit-test: deadcode-test
	go test -race $(shell go list ./... | grep -v tests)

.PHONY: e2e-test
e2e-test: deadcode-test
	go test -tags e2e -race ./tests/...

# to find uninitialized dependencies
.PHONY: deadcode-test
deadcode-test:
	gometalinter --vendor --disable-all --enable=deadcode di

.PHONY: vet
vet:
	@gometalinter --vendor --disable-all --enable=vet ./...

.PHONY: deadcode
deadcode:
	@gometalinter --vendor --disable-all --enable=deadcode ./...

.PHONY: lint
lint:
	# ignore comments for exported objects
	# ignore Err prefix
	gometalinter --vendor --disable-all --enable=golint --exclude="(should have comment|ErrFoo)" ./...

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out -race $(shell go list ./... | grep -v tests)
	go tool cover -html=coverage.out

.PHONY: brew-update
brew-update:
	bash .circleci/scripts/entrypoint.bash $(VERSION)
