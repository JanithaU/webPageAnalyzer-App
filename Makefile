# Variables
BINARY_NAME=webAnalyzer

# Build the application
build: tidy
	echo "Building..."
	go build -o $(BINARY_NAME) cmd/web/main.go

# Run the application
run: build
	./$(BINARY_NAME)

# Clean up generated files
clean:
	go clean
	rm -f $(BINARY_NAME)

# Tidy up module dependencies
tidy:
	echo "downloading dependencies.."
	go mod tidy

test:
	echo "Testing..."
	go test ./... -cover