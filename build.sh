#!/bin/bash

# This script is used to build and run the server locally using Docker for 
# testing purposes prior to pushing to Azure Container Registry

# Load environment variables from the .env file
source .env

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
    --env APP_SECRET="$APP_SECRET" \
    go-web-starter:$version