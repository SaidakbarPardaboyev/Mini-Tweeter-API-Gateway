FROM golang:1.24-alpine AS builder

# Install necessary dependencies for building the application and Cgo support
RUN apk update && apk add --no-cache \
    build-base \
    libwebp-dev \
    curl

WORKDIR /go/src/github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway

# Copy the local package files to the container's workspace.
COPY . ./

# Make sure to copy the config folder into the builder container
COPY ./config /go/src/github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config

# Install dependencies and build the application
RUN export CGO_ENABLED=1 && \
    export GOOS=linux && \
    go mod vendor && \
    go build -o /Mini_Tweeter_API_Gateway ./cmd/main.go

# Use a minimal Alpine image for the final runtime stage
FROM alpine

# Copy the built Go binary from the builder stage
COPY --from=builder /Mini_Tweeter_API_Gateway .

# Copy the config directory from the builder stage to the final image
COPY --from=builder /go/src/github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config /config

# Install necessary runtime dependencies (e.g., curl)
RUN apk add --no-cache curl

# Set the entry point for the container
ENTRYPOINT ["/Mini_Tweeter_API_Gateway"]
