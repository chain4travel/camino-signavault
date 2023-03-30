#!/bin/bash

#
# This script will build the signavault binary in the project's root directory.
# The script must be executed from the root of the project.
#

OUTPUT="signavault"

echo "Building signavault..."
go build -o "$OUTPUT" cmd/camino-signavault/main.go
echo "SUCCESS"