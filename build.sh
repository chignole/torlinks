#!/bin/bash

# Get the current date and time in the desired format
build_time=$(date +%Y%m%d-%H%M)

# Compile the application, injecting the build date and time
go build -ldflags "-X chignole/torlinks/cmd.buildDate=$build_time" -o build/torlinks main.go
