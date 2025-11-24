# ---- config ----
BIN := 412alloc          # The executable name the autograder expects
PKG ?= .                  # Go package to build (default for this is current directory)
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
  --exclude='$(BIN)' \
  --exclude='test_files/correct' \
  --exclude='bs81.log' \
  --exclude='CodeCheck2' \

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
	./$(BIN) -s test_files/ex1.txt

default-parse:
	./$(BIN) -r test_files/ex1.txt

wrong-parse:
	./$(BIN) -r test_files/Holmes.txt

tar: clean
	@tar $(TAR_EXCLUDES) -cvf $(TARFILE) .

# commands for lab 2 

default-rename: clean-build
	./$(BIN) -x test_files/rename.txt

ccOne: clean-build
	./$(BIN) -x test_files/cc1.txt

ccTwo: clean-build
	./$(BIN) -x test_files/cc2.i

ccThree: clean-build
	./$(BIN) -x test_files/cc3.i

ccFour: clean-build
	./$(BIN) -x test_files/cc4.i

ccFive: clean-build
	./$(BIN) -x test_files/cci.i

# checkin 2

default-alloc: clean-build
	./$(BIN) 5 test_files/rename.txt

cc2One: clean-build
	./$(BIN) 5 test_files/cc1.txt

cc2One2: clean-build
	./$(BIN) 7 test_files/cc1.txt

cc2Two: clean-build
	./$(BIN) 5 test_files/cc2.i

cc2Two2: clean-build
	./$(BIN) 7 test_files/cc2.i

cc2Three: clean-build
	./$(BIN) 5 test_files/cc3.i

cc2Four: clean-build
	./$(BIN) 5 test_files/cc4.i

cc2Five: clean-build
	./$(BIN) 5 test_files/cc5.i

# autograder

r2: clean-build
	./$(BIN) 3 test_files/report2.txt

r2rename: clean-build
	./$(BIN) -x test_files/report2.txt


ag2: clean-build
	./$(BIN) 3 test_files/report2.i

# lab 3
easySched: clean-build
	./$(BIN) test_files/easySched.txt


