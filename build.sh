#!/bin/bash

# Set environment variables
export DB_SERVER="SQL URL"
export DB_PORT="SQL PORT"
export DB_USER="SQL USER"
export DB_PASSWORD="SQL PASS"
export DB_NAME="SQL DB"

# Read the version from the command line argument or default to "latest"
version="${1:-latest}"

# Build docker image for deployment to Azure
docker buildx build \
    --build-arg VERSION="$version" \
    -f Dockerfile \
    --platform linux/amd64 \
    --tag "go-web-starter:$version" .

# Run the docker image locally
docker run \
    -p 8080:8080 \
    --env DB_SERVER="$DB_SERVER" \
    --env DB_PORT="$DB_PORT" \
    --env DB_USER="$DB_USER" \
    --env DB_PASSWORD="$DB_PASSWORD" \
    --env DB_NAME="$DB_NAME" \
    go-web-starter:$version