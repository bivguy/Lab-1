# ---- config ----
BIN := 412fe          # The executable name the autograder expects
PKG ?= .              # Go package to build (default for this is current directory)
NETID ?= bs81
TARFILE ?= ../$(NETID).tar

# Exclusions 
TAR_EXCLUDES := \
  --exclude='*_test.go' \
  --exclude='testdata' \
  --exclude='_tests' \
  --exclude='test_files' \
  --exclude='scanner/scanner_tests' \
  --exclude='parser/parser_tests' \
  --exclude='l1auto' \
  --exclude='*.tar' \
  --exclude='Grader' \
  --exclude='Timer' \
  --exclude='*.log' \
  --exclude='.git' \
  --exclude='.vscode' \
  --exclude='$(BIN)'

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

default-parse:
	./412fe -r test_files/ex1.txt

wrong-parse:
	./412fe -r test_files/Holmes.txt


tar: clean
	@tar $(TAR_EXCLUDES) -cvf $(TARFILE) .



