# Stage 1: Build the Go application
FROM golang:1.23.2-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules dependencies
RUN go mod tidy

# Install goose tool
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy the entire project to the container
COPY . .

# Set CGO_ENABLED=0 to build a statically-linked binary
# GOOS=linux and GOARCH=amd64 to ensure the binary is built for Linux architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /build/billing-loan-system .

# Stage 2: Create the minimal runtime image
FROM alpine:latest  

# Install required dependencies (if any)
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /build/billing-loan-system /usr/local/bin/billing-loan-system

# Copy the configuration file
COPY --from=builder /build/billing-loan-system .
COPY ./configs/env.yml /root/configs/
COPY ./configs/env.yml /root/configs/env-local.yml
RUN ls -a

# Expose the port your application will run on (e.g., HTTP server on 8088)
EXPOSE 8088

# Command to run the binary
CMD ["billing-loan-system", "http"]
