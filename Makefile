ifndef VERBOSE
.SILENT:
endif

all: binary

.PHONY: fmt
fmt:
	sh ./scripts/build.sh fmt

.PHONY: clean
clean:
	sh ./scripts/build.sh clean

.PHONY: binary
binary:
	sh ./scripts/build.sh build

.PHONY: install
install:
	sh ./scripts/build.sh install