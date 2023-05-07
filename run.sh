#!/bin/bash

# Set environment variables
export DB_SERVER="SQL URL"
export DB_PORT="SQL PORT"
export DB_USER="SQL USER"
export DB_PASSWORD="SQL PASS"
export DB_NAME="SQL DB"

# Run the Go application
go run -ldflags "-X main.version=latest" ./src/

