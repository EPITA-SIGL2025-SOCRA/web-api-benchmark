FROM golang:1.20-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Copy local code to the container image.
COPY go/web-service/ ./

# Build the binary.
RUN go build -v -o server

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY ./data/tractors.json /data/tractors.json

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

ENV DATA_JSON_FILE_PATH=/data/tractors.json
EXPOSE 8080
# Run the web service on container startup.
CMD ["/app/server"]
