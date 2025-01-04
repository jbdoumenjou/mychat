# Step 1: Build the executable binary
# Use the official Golang image to create a build artifact.
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /src

# Download dependencies as a separate step to take advantage of Docker's caching.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

# Set the target platform for cross-compilation
ARG TARGETOS
ARG TARGETARCH

# Build project using bind mounts to avoid having to copy everything into the container.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -installsuffix cgo -o /bin/mychat ./cmd/main.go

# Step 2: Use a Docker multi-stage build to create a lean production image
# Start from scratch for a smaller final image
FROM alpine:3.20

#  Get the last certificates & tz-data version.
RUN apk --no-cache --no-progress add ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

RUN adduser \
    --disabled-password \
    --home /dev/null \
    --no-create-home \
    --shell /sbin/nologin \
    --gecos mychatuser \
    --uid 10000 \
    mychatuser

USER mychatuser

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /bin/mychat /bin

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/bin/mychat"]
