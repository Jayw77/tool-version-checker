# Use the official Golang image as a parent image
FROM golang:alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download all the dependencies
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o version-checker .

# Use a Docker multi-stage build to create a lean production image
# Start with a scratch image
FROM alpine:latest

# Copy the compiled application from the builder stage
COPY --from=builder /app/version-checker .

# Copy the relevant files from the builder stage
COPY --from=builder /app/home.html .
COPY --from=builder /app/default_config.yaml .

# Install tail or any other utilities you might need
RUN apk add --no-cache bash

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./version-checker"]
