#!/bin/bash

# Stop the script on any error
set -e

# Function to create ubuntu-rootfs if it doesn't already exist
function ensure_ubuntu_rootfs() {
    if [ ! -d "ubuntu-rootfs" ]; then
        echo "Creating ubuntu-rootfs..."
        if ! sudo debootstrap focal ubuntu-rootfs; then
            echo "Failed to create ubuntu-rootfs. Trying to install debootstrap..."
            return 1
        fi
    else
        echo "ubuntu-rootfs already exists."
    fi
    return 0
}

# Function to check the existence of debootstrap and install if needed
function ensure_debootstrap() {
    if ! command -v debootstrap &> /dev/null; then
        echo "debootstrap not found, installing..."
        sudo apt-get update
        sudo apt-get install -y debootstrap
    else
        echo "debootstrap is already installed."
    fi
}

# Ensure rootfs is created, installing debootstrap if necessary
echo "Checking system setup for building containers..."
if ! ensure_ubuntu_rootfs; then
    ensure_debootstrap
    # Retry creating ubuntu-rootfs after installing debootstrap
    ensure_ubuntu_rootfs
fi

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o docklet
if [ $? -ne 0 ]; then
    echo "Failed to build for Linux"
    exit 1
fi

# # Build for Windows
# echo "Building for Windows..."
# GOOS=windows GOARCH=amd64 go build -o docklet.exe
# if [ $? -ne 0 ]; then
#     echo "Failed to build for Windows"
#     exit 1
# fi

echo "Builds completed successfully!"
