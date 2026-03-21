FROM golang:1.25-alpine AS builder

# Set the working directory outside the GOPATH to enable the support for modules.
WORKDIR /app

# Install git. Run 'docker build --no-cache .' to update dependencies.
# Also install ca-certificates and tzdata which are useful.
RUN apk add --no-cache git ca-certificates tzdata

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build the executable to `/app/bot`.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bot ./cmd/bot

# Use a minimal Alpine distribution as the runner image.
FROM alpine:3.19

# Useful if we are doing HTTPS requests.
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the compiled binary from the builder context.
COPY --from=builder /app/bot .

# Command to run the application
CMD ["./bot"]
