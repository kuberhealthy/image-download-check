IMAGE := "kuberhealthy/image-download-check"
TAG := "latest"

# Build the image download check container locally.
build:
	podman build -f Containerfile -t {{IMAGE}}:{{TAG}} .

# Run the unit tests for the image download check.
test:
	go test ./...

# Build the image download check binary locally.
binary:
	go build -o bin/image-download-check ./cmd/image-download-check
