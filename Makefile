TAG?=$(shell git describe --abbrev=0 --tags 2>/dev/null || echo "v0.0.0" )
COMMIT?=$(shell git rev-parse HEAD)

default: build

.PHONY: build
build:
	COMMIT=$(COMMIT) TAG=$(TAG) ./scripts/build.sh

.PHONY: test
test:
	./scripts/test.sh

.PHONY: clean
clean:
	./scripts/clean.sh

.PHONY: help
help:
	@echo "Usage:"
	@echo "	make build		build binary files"
	@echo "	make test		run unit tests"
	@echo "	make clean		clean up built files"
	@echo "	make help		show this message"
