# ---- config ----
BIN := 412fe          # The executable name the autograder expects
PKG ?= .              # Go package to build (default: current directory)

# ---- targets ----
.PHONY: all build clean run test
all: build

build:
	@echo ">> building $(BIN) from $(PKG)"
	go build -o $(BIN) $(PKG)

clean:
	@echo ">> cleaning"
	@rm -f $(BIN)

clean-build:
	@rm -f $(BIN)
	go build -o $(BIN) $(PKG)

default-scan:
	./412fe -s test_files/ex1.txt



