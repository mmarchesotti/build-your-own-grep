# Makefile for the mygrep project

# --- Variables ---
# The name of the final binary
BINARY_NAME=mygrep
# The path to the main package
CMD_PATH=./cmd/mygrep

# --- Targets ---

# The .PHONY directive tells make that these are not files.
# This ensures the command runs even if a file with the same name exists.
.PHONY: all build run test clean help

# 'all' is a standard default target. We'll have it run the build.
all: build

# build: Compiles the Go source code into a binary.
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "$(BINARY_NAME) built successfully."

# run: Builds the project first, then runs the binary.
# You can pass arguments to the program using the ARGS variable.
# Example: make run ARGS="-E 'appl.*' fruits.txt"
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME) $(ARGS)

# test: Runs all tests in the project.
test:
	@echo "Running tests..."
	@go test ./...

# clean: Removes the compiled binary file.
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# help: Prints this help message.
help:
	@echo "Available commands:"
	@echo "  make build    - Compiles the application."
	@echo "  make run      - Builds and runs the application. Use ARGS to pass arguments."
	@echo "  make test     - Runs all tests."
	@echo "  make clean    - Removes the built binary."
	@echo "  make help     - Shows this help message."
