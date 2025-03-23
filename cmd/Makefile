# set environment variables for cross-compilation
GOOS=linux
GOARCH=amd64
BINARY_NAME=health-lambda
BUILD_DIR=cmd/health-monitor

.PHONY: build clean package deploy test build-HealthMonitorFunction

build: # build the Go binary for AWS Lambda
	@echo "Building Lambda binary..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/main.go

clean: # remove compiled binaries
	@echo "Cleaning up..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

test: # run Go tests
	@echo "Running tests..."
	go test ./...

package: # package the application with AWS SAM
	@echo "Packaging application..."
	sam build

deploy: # deploy using AWS SAM
	@echo "Deploying application..."
	sam deploy --guided

build-HealthMonitorFunction: # special target for SAM builds
	@echo "Building for SAM deployment..."
	GOOS=linux GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/bootstrap ./health-monitor/main.go

