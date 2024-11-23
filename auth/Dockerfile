FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

COPY . /app

# Build the binary.
RUN go build -ldflags="-s -w" -v -o main ./cmd/api

FROM alpine:latest AS runtime

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/main /

EXPOSE 8081

CMD ["./main"]
