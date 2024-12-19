PROGRAM_NAME = ghligh
BIN_DIR = $(HOME)/go/bin
SRC = $(wildcard *.go)
COMPLETION_DIR := $(HOME)/.local/share/bash-completion/completions/

COMPLETION_SCRIPT = ghligh.bash

.PHONY: all
all: build

.PHONY: run
run:
	go run $(SRC)

.PHONY: build
build:
	go build -o $(PROGRAM_NAME)

.PHONY: install
install:
	@echo "installing $(PROGRAM_NAME)..."
	@echo "to use it add $(BIN_DIR) to your PATH"
	go install
	@echo "installing completion file"
	mkdir -p $(COMPLETION_DIR)
	install -m 0644 $(COMPLETION_SCRIPT) $(COMPLETION_DIR)

.PHONY: uninstall
uninstall:
	@echo "uninstalling $(PROGRAM_NAME)..."
	@echo "deleting file $(PROGRAM_NAME) from ~/go/bin/"
	[ -f "$(BIN_DIR)/$(PROGRAM_NAME)" ] && rm "$(BIN_DIR)/$(PROGRAM_NAME)"
	@echo "removing completion file"
	[ -f "$(COMPLETION_DIR)/$(COMPLETION_SCRIPT)" ] && rm "$(COMPLETION_DIR)/$(COMPLETION_SCRIPT)"




.PHONY: clean
clean:
	rm -f $(PROGRAM_NAME)
