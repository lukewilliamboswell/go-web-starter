#!/bin/bash

# This script is used to run the server locally for testing purposes

# Load environment variables from the .env file
source .env

# Run the Go application without API authentication
go run -ldflags "-X main.version=dev" ./src/

