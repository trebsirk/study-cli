#!/bin/bash

# Set the module name
MODULE_NAME="github.com/trebsirk/study-cli"

# Initialize the Go module
echo "Initializing Go module..."
go mod init $MODULE_NAME

# Install Cobra
echo "Installing Cobra..."
go get -u github.com/spf13/cobra@latest

# Install JSON package
echo "Installing JSON package..."
go get github.com/lib/pq

# for colorful text output
go get github.com/fatih/color

source setup.sh