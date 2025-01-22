#!/bin/sh

# Copyright 2025 Defense Unicorns
# SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

# Define variables
# renovate: datasource=docker depName=ghcr.io/defenseunicorns/uds-security-hub versioning=loose
TAG="v2.2.0-20250107.0413"
IMAGE="ghcr.io/defenseunicorns/uds-security-hub:$TAG"
OUTPUT_PATH="./uds-security-hub.db"

mkdir -p artifacts
METADATA_OUTPUT_PATH="./artifacts/security-hub-metadata.json"

# Pull and run image
echo "Pulling Docker image: $IMAGE"
docker pull "$IMAGE"
echo "Running Docker container from image: $IMAGE"
CONTAINER_ID=$(docker run -d "$IMAGE" tail -f /dev/null)

# Check if the container started successfully
if [ -z "$CONTAINER_ID" ]; then
    echo "Failed to start container."
    exit 1
fi
echo "Container started with ID: $CONTAINER_ID"

# Copy the file from the container
echo "Copying file from container: $CONTAINER_ID"
docker cp "$CONTAINER_ID:/data/uds/uds-security-hub.db" "$OUTPUT_PATH"
docker cp "$CONTAINER_ID:/data/uds/artifacts/security-hub-metadata.json" "$METADATA_OUTPUT_PATH"

# Check if the file was copied successfully
if [ $? -eq 0 ]; then
    echo "File copied successfully to $OUTPUT_PATH"
else
    echo "Failed to copy file from container."
    # Stop and remove the container before exiting
    docker rm -f "$CONTAINER_ID"
    exit 1
fi

# Stop and remove the container
echo "Stopping and removing container: $CONTAINER_ID"
docker rm -f "$CONTAINER_ID"

echo "UDS Security Hub database downloaded successfully."
