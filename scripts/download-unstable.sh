#!/bin/bash

# Set the image name and tag
image_name="quay.io/testnetworkfunction/cnf-certification-test"
image_tag="unstable"

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
	echo "Failed to pull Docker image ${image_name}:${image_tag}"
fi
