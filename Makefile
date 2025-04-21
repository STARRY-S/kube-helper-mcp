TAG?=$(shell git describe --abbrev=0 --tags 2>/dev/null || echo "v0.0.0")
COMMIT?=$(shell git rev-parse HEAD)

default: build

.PHONY: prepare
prepare:
	uv venv
	uv pip install mcpo

.PHONY: build
build:
	COMMIT=$(COMMIT) TAG=$(TAG) goreleaser build --snapshot --clean

.PHONY: test
test:
	./scripts/test.sh

.PHONY: serve
serve:
	KUBE_CONFIG=${HOME}/.kube/config uvx mcpo --config ./mcpo/config.json

.PHONY: clean
clean:
	./scripts/clean.sh

.PHONY: help
help:
	@echo "Usage:"
	@echo "	make prepare    init the pytho venv and install dependencies"
	@echo "	make serve		run the mcpo server"
	@echo "	make build		build binary files"
	@echo "	make test		run unit tests"
	@echo "	make clean		clean up built files"
	@echo "	make help		show this message"
