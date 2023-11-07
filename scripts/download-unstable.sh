#!/bin/bash

# Check for LOCAL_TESTING environment variable and exit early
if [[ -z "$FORCE_DOWNLOAD_UNSTABLE" ]]; then
	echo "Skipping download of unstable image"
	exit 0
fi

# If docker isn't installed, exit gracefully
if ! command -v docker &>/dev/null; then
	echo >&2 "Docker is not installed. Skipping download of unstable image"
	exit 0
fi

# Set the image name and tag
image_name=quay.io/testnetworkfunction/cnf-certification-test
image_tag=unstable

# Delete the image if it exists
if docker image inspect "${image_name}:${image_tag}" >/dev/null; then
	docker rmi "${image_name}:${image_tag}"
fi

# Pull the image
docker pull "${image_name}:${image_tag}"

# Check if the image was pulled successfully
if docker image inspect "${image_name}:${image_tag}" >/dev/null; then
	echo "Docker image ${image_name}:${image_tag} was pulled successfully"
else
	echo >&2 "Failed to pull Docker image ${image_name}:${image_tag}"
	exit 1
fi
