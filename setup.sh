#!/bin/bash

# for source .env to work, run source setup.sh, not ./setup.sh

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found"
    exit 1
fi

# Run the source command to load the environment variables
source .env