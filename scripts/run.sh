#!/bin/bash

echo "Starting SSH AI Server..."
echo "Server will listen on port 2212"
echo "Connect with: ssh gpt-5@localhost -p 2212"
echo ""

go run main.go