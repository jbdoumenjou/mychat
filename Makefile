# Define the base name for the image
IMAGE_BASE_NAME=mychat

# Define the container name
CONTAINER_NAME=mychat

# Retrieve the current Git commit hash
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)

# Combine base name with commit hash for the final image tag
IMAGE_TAG=$(IMAGE_BASE_NAME):$(GIT_COMMIT_HASH)

# Default build target
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

.PHONY: help test lint build build-image run-image clean-image up down restart gen-doc

help: ## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST) | column -tl 2

test: ## Test go files and report coverage.
	go test -v -race -cover ./...

lint: ## List all the linting issues.
	golangci-lint run

build:  ## Build the application.
	@echo "Building for OS: $(GOOS), Arch: $(GOARCH)"
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo -o mychat ./cmd/main.go

build-image: ## Build the Docker image.
	@echo "Building Docker image with tag ${IMAGE_TAG}..."
	docker build -t $(IMAGE_TAG) .

run-image: build-image ## Run the Docker container.
	@echo "Running Docker container..."
	docker run -d -p 8080:8080 --name $(CONTAINER_NAME) $(IMAGE_TAG)

clean-image: ## Stop and remove the Docker container.
	@echo "Cleaning up Docker container..."
	-docker stop $(CONTAINER_NAME)
	-docker rm $(CONTAINER_NAME)
	-docker rmi $(IMAGE_TAG)

up: ## Start the docker-compose with the database and the app.
	docker-compose up --build -d

down: ## Stop the docker-compose with the database and the app.
	docker-compose down

restart: down up ## Restart the docker-compose with the database and the app.
